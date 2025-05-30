#ifndef CT_C24_LW_LIBRARIES_CEPEXAAA_OPENINGDATA_H
#define CT_C24_LW_LIBRARIES_CEPEXAAA_OPENINGDATA_H
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <fftw3.h>
#include <libavutil/file.h>
#include <stdio.h>
#include <libswresample/swresample.h>
#include <libavutil/opt.h>
double *process(char *arg, int *return_size, const int CHANNEL_NUMBER, int *sample_rate, int streamIndex, AVCodecContext *codec_context, AVFormatContext *context);
int openFile(char *arg, AVCodecContext **codec_context, const int *sample_rate, AVFormatContext **context,  int *streamIndex);
#endif //CT_C24_LW_LIBRARIES_CEPEXAAA_OPENINGDATA_H
