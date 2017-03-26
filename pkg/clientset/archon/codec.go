/*
Copyright 2017 The Archon Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package archon

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/pkg/api"
)

// DirectDefaultingCodecFactory provides DirectDefaultDecoder
type DirectDefaultingCodecFactory struct {
	serializer.CodecFactory
}

// DecoderToVersion returns an decoder that does not do conversion. gv is ignored.
func (f DirectDefaultingCodecFactory) DecoderToVersion(serializer runtime.Decoder, _ runtime.GroupVersioner) runtime.Decoder {
	return DirectDefaultingDecoder{
		Decoder: serializer,
	}
}

type DirectDefaultingDecoder struct {
	runtime.Decoder
}

// Decode does not do conversion. It removes the gvk during deserialization. Plus it sets
// defaults on obj.
func (d DirectDefaultingDecoder) Decode(data []byte, defaults *schema.GroupVersionKind, into runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {
	obj, gvk, err := d.Decoder.Decode(data, defaults, into)
	if obj != nil {
		kind := obj.GetObjectKind()
		// clearing the gvk is just a convention of a codec
		kind.SetGroupVersionKind(schema.GroupVersionKind{})
		api.Scheme.Default(obj)
	}
	return obj, gvk, err
}
