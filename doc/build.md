
# Package Config

libwhisper.pc

```pkg-config
prefix=/Users/djt/Projects/go-whisper/

Name: libwhisper
Description: Whisper is a C/C++ library for speech transcription, translation and diarization.
Version: 0.0.0
Cflags: -I${prefix}/third_party/whisper.cpp/include -I${prefix}/third_party/whisper.cpp/ggml/include
Libs: -L${prefix}/third_party/whisper.cpp -lwhisper -lggml -lm -lstdc++
```

libwhisper-linux.pc

```pkg-config
prefix=/Users/djt/Projects/go-whisper/

Name: libwhisper-linux
Description: Whisper is a C/C++ library for speech transcription, translation and diarization.
Version: 0.0.0
```

libwhisper-darwin.pc

```pkg-config
prefix=/Users/djt/Projects/go-whisper/

Name: libwhisper-darwin
Description: Whisper is a C/C++ library for speech transcription, translation and diarization.
Version: 0.0.0
Libs: -framework Accelerate -framework Metal -framework Foundation -framework CoreGraphics
```

I don't know what the windows one should be as I don't have a windows machine.
