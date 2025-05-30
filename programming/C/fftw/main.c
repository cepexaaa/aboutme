#include "checkException.h"
#include "crossCorrelation.h"
#include "openingData.h"
#include "return_codes.h"

#include <fftw3.h>
#include <stdio.h>
#define MAX(a, b) (a > b ? b : a)

int main(int argc, char *argv[])
{
	AVFormatContext *context1 = NULL, *context2 = NULL;
	AVCodecContext *codec_context1 = NULL, *codec_context2 = NULL;
	double *freq1 = NULL, *freq2 = NULL;
	int size1 = 0, size2 = 0, sample_rate1 = -1, sample_rate2 = -1, streamIndex1, streamIndex2;
	int shift = 0;
	if (argc < 2)
	{
		return ERROR_ARGUMENTS_INVALID;
	}
	int extension_file1 = check_file_extension(argv[1]);
	if (extension_file1 == -1)
	{
		return ERROR_DATA_INVALID;
	}
	int sample_rate = -1;
	if (argc == 3)
	{
		int extension_file2 = check_file_extension(argv[2]);
		if (extension_file1 != extension_file2)
		{
			return ERROR_ARGUMENTS_INVALID;
		}
		int exception1 = openFile(argv[1], &codec_context1, &sample_rate1, &context1, &streamIndex1);
		if (exception1)
		{
			fprintf(stderr, "uncorrect data\n");
			return ERROR_DATA_INVALID;
		}
		int exception2 = openFile(argv[2], &codec_context2, &sample_rate2, &context2, &streamIndex2);
		if (exception2)
		{
			fprintf(stderr, "uncorrect data\n");
			return ERROR_DATA_INVALID;
		}
		sample_rate = MAX(sample_rate1, sample_rate2);
		freq1 = process(argv[1], &size1, 0, &sample_rate, streamIndex1, codec_context1, context1);
		freq2 = process(argv[2], &size2, 0, &sample_rate, streamIndex2, codec_context2, context2);
		if (freq1 == NULL || freq2 == NULL)
		{
			fprintf(stderr, "uncorrect data\n");
			return ERROR_DATA_INVALID;
		}
		shift = crossing(freq1, freq2, &size1, &size2);
	}
	else
	{
		int exception1 = openFile(argv[1], &codec_context1, &sample_rate1, &context1, &streamIndex1);
		sample_rate = sample_rate1;
		freq1 = process(argv[1], &size1, 0, &sample_rate1, streamIndex1, codec_context1, context1);
		av_free(freq1);
	}
	if (shift == -1)
	{
		fprintf(stderr, "Data is wrong\n");
		return ERROR_ARGUMENTS_INVALID;
	}
	printf("delta: %i samples\nsample rate: %i Hz\ndelta time: %i ms\n", shift, sample_rate, (int)(((long)shift * 1000) / sample_rate));
	return SUCCESS;
}
