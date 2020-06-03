/*
 * Copyright 2020 Mandelsoft. All rights reserved.
 *  This file is licensed under the Apache Software License, v. 2 except as noted
 *  otherwise in the LICENSE file
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package kipxe

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gardener/controller-manager-library/pkg/logger"
)

func Hash(key string) string {
	data := sha1.Sum([]byte(key))
	return hex.EncodeToString(data[:])
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

////////////////////////////////////////////////////////////////////////////////

type Cache interface {
	Bytes(url *url.URL) ([]byte, error)
	Serve(url *url.URL, w http.ResponseWriter, r *http.Request)
	Cleanup(logger logger.LogContext, ttl time.Duration)
}

type DirCache struct {
	lock sync.Mutex
	logger.LogContext
	path    string
	actions map[string]*cacheAction
}

////////////////////////////////////////////////////////////////////////////////

type CacheAction struct {
	lock sync.Mutex
	ref  *cacheAction
}

func (this *CacheAction) Execute(f func()) {
	this.lock.Lock()
	defer this.lock.Unlock()
	f()
}

func (this *CacheAction) Bytes() ([]byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.ref == nil {
		return nil, fmt.Errorf("outdated")
	}
	return this.ref.Bytes()
}

func (this *CacheAction) Serve(w http.ResponseWriter, r *http.Request) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.ref == nil {
		return
	}
	this.ref.Serve(w, r)
}

func (this *CacheAction) Done() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.ref != nil {
		this.ref.release()
		this.ref = nil
	}
}

////////////////////////////////////////////////////////////////////////////////

type cacheAction struct {
	sync.RWMutex
	cache    *DirCache
	usecount int
	url      *url.URL
	key      string
	base     string
}

func (this *cacheAction) release() {
	this.Lock()
	defer this.Unlock()

	this.usecount--
	if this.usecount <= 0 {
		this.cache.release(this.key)
	}
}

func write(w io.Writer, data []byte) error {
	start := 0
	for start < len(data) {
		n, err := w.Write(data[start:])
		if n > 0 {
			start += n
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *cacheAction) fill(writer io.Writer) error {
	err := this._fill(writer)
	if err != nil {
		this.cache.remove(this.base)
	}
	return err
}

func (this *cacheAction) _fill(writer io.Writer) error {
	this.cache.Infof("caching %s [%s]", this.url, this.base)
	file, err := os.OpenFile(this.base, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(this.url.String())
	if err != nil {
		return fmt.Errorf("URL get failed: %s", err)
	}

	defer resp.Body.Close()
	meta := ResourceMeta{}
	meta[CONTENT_URL] = this.url.String()
	if t := resp.Header.Get(CONTENT_TYPE); t != "" {
		meta[CONTENT_TYPE] = t
	}
	meta.Write(this.base)

	var tmp [8196]byte
	var fail error
	var wfail error
	for {
		n, err := resp.Body.Read(tmp[:])
		if n > 0 {
			if fail == nil {
				fail = write(file, tmp[:n])
				if writer == nil {
					return fail
				}
			}
			if wfail == nil && writer != nil {
				wfail = write(writer, tmp[:n])
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n < 0 {
			break
		}
	}
	return nil
}

func (this *cacheAction) _bytes() ([]byte, error) {
	if fileExists(this.base) {
		buf := &bytes.Buffer{}
		file, err := os.OpenFile(this.base, os.O_RDONLY, 0660)
		if err != nil {
			this.cache.remove(this.base)
			return nil, err
		}
		defer file.Close()
		var tmp [8096]byte
		for {
			n, err := file.Read(tmp[:])
			if n > 0 {
				buf.Write(tmp[:n])
			}
			if err != nil {
				if err == io.EOF {
					return buf.Bytes(), nil
				}
				return nil, err
			}
			if n < 0 {
				return buf.Bytes(), nil
			}
		}
	}
	return nil, nil
}

func (this *cacheAction) bytes() ([]byte, error) {
	this.RLock()
	defer this.RUnlock()
	return this._bytes()
}

func (this *cacheAction) Bytes() ([]byte, error) {

	data, err := this.bytes()
	if data != nil || err != nil {
		return data, err
	}

	this.Lock()
	defer this.Unlock()

	data, err = this._bytes()
	if data != nil || err != nil {
		return data, err
	}

	buf := &bytes.Buffer{}
	err = this.fill(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (this *cacheAction) _serve(w http.ResponseWriter, r *http.Request) bool {
	if fileExists(this.base) {
		meta := ResourceMeta{}
		meta.Read(this.base)
		if meta[CONTENT_TYPE] != "" {
			w.Header().Set(CONTENT_TYPE, meta[CONTENT_TYPE])
		}
		http.ServeFile(w, r, this.base)
		return true
	}
	return false
}

func (this *cacheAction) serve(w http.ResponseWriter, r *http.Request) bool {
	this.RLock()
	defer this.RUnlock()
	return this._serve(w, r)
}

func (this *cacheAction) Serve(w http.ResponseWriter, r *http.Request) {
	if this.serve(w, r) {
		return
	}

	this.Lock()
	defer this.Unlock()

	if this._serve(w, r) {
		return
	}

	this.fill(w)
}

////////////////////////////////////////////////////////////////////////////////

type ResourceMeta map[string]string

func (this ResourceMeta) Write(base string) {
	file, err := os.OpenFile(cachemeta(base), os.O_WRONLY|os.O_CREATE, 0660)
	if err == nil {
		defer file.Close()
		for k, v := range this {
			file.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}
}

func (this ResourceMeta) Read(base string) {
	file, err := os.OpenFile(cachemeta(base), os.O_RDONLY, 0660)
	if err == nil {
		defer file.Close()
		r := bufio.NewReader(file)
		for {
			line, err := r.ReadString('\n')
			if !strings.HasPrefix(line, "#") {
				i := strings.Index(line, ": ")
				if i > 0 {
					this[line[:i]] = line[i+2:]
				}
			}
			if err != nil {
				break
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

func cachemeta(base string) string {
	return base + ".meta"
}
func iscachemeta(base string) bool {
	return strings.HasPrefix(base, ".meta")
}

func NewDirectoryCache(logger logger.LogContext, path string) (*DirCache, error) {
	err := os.MkdirAll(path, 0770)
	if err != nil {
		return nil, err
	}
	return &DirCache{
		LogContext: logger.NewContext("server", "cache"),
		path:       path,
		actions:    map[string]*cacheAction{},
	}, nil
}
func (this *DirCache) remove(base string) {
	os.Remove(base)
	os.Remove(cachemeta(base))
}

func (this *DirCache) release(key string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.actions, key)
}

func (this *DirCache) GetAction(url *url.URL) *CacheAction {
	return this.getAction(Hash(url.String()), url)
}

func (this *DirCache) GetActionForKey(key string) *CacheAction {
	return this.getAction(key, nil)
}

func (this *DirCache) getAction(key string, url *url.URL) *CacheAction {
	base := filepath.Join(this.path, key)

	this.lock.Lock()
	defer this.lock.Unlock()

	action := this.actions[key]
	if action == nil {
		action = &cacheAction{
			url:   url,
			key:   key,
			base:  base,
			cache: this,
		}
		this.actions[key] = action
	} else {
		if url != nil {
			action.url = url
		}
	}
	action.usecount++
	return &CacheAction{ref: action}
}

func (this *DirCache) Bytes(url *url.URL) ([]byte, error) {
	action := this.GetAction(url)
	defer action.Done()
	return action.Bytes()
}

func (this *DirCache) Serve(url *url.URL, w http.ResponseWriter, r *http.Request) {
	action := this.GetAction(url)
	defer action.Done()
	action.Serve(w, r)
}

func (this *DirCache) Cleanup(logger logger.LogContext, duration time.Duration) {
	if logger == nil {
		logger = this
	}
	files, err := ioutil.ReadDir(this.path)
	if err != nil {
		return
	}

	now := time.Now()
	for _, f := range files {
		action := this.GetActionForKey(f.Name())
		action.Execute(func() {
			if !iscachemeta(f.Name()) {
				fpath := filepath.Join(this.path, f.Name())
				finfo, err := os.Stat(fpath)
				if err == nil {
					ts := finfo.Sys().(*syscall.Stat_t).Atim
					t := time.Unix(int64(ts.Sec), int64(ts.Nsec))

					now.Sub(t)
					if now.Sub(t) > duration {
						logger.Infof("cleanup %s [%s]", f.Name(), now.Sub(t))
						this.remove(fpath)
					}
				}
			}
		})
	}
}
