#ifndef CT_C24_LW_LIBRARIES_CEPEXAAA_CROSSCORRELATION_H
#define CT_C24_LW_LIBRARIES_CEPEXAAA_CROSSCORRELATION_H
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <fftw3.h>
#include <libavutil/file.h>
int crossing(double *freq1, double *freq2, int *size1, int *size2);
#endif //CT_C24_LW_LIBRARIES_CEPEXAAA_CROSSCORRELATION_H
