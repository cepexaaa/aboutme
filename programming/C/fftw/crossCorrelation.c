#include <libavformat/avformat.h>

#include <fftw3.h>
#define MIN(a, b) (a > b ? b : a)
int crossing(double *freq1, double *freq2, int *size1, int *size2)
{
	long bothsize = *size2 + *size1 - 1;
	int max_index = -1, max_value = INT_MIN;
	freq1 = (double *)av_realloc(freq1, bothsize * sizeof(double));
	freq2 = (double *)av_realloc(freq2, bothsize * sizeof(double));
	for (int i = *size1; i < bothsize; i++)
	{
		freq1[i] = 0.0;
	}
	for (int i = *size2; i < bothsize; i++)
	{
		freq2[i] = 0.0;
	}
	fftw_complex *freq_signal1 = (fftw_complex *)fftw_malloc(sizeof(fftw_complex) * bothsize);
	fftw_plan plan1 = fftw_plan_dft_r2c_1d(bothsize, freq1, freq_signal1, FFTW_ESTIMATE);
	fftw_execute(plan1);
	fftw_complex *freq_signal2 = (fftw_complex *)fftw_malloc(sizeof(fftw_complex) * bothsize);
	fftw_plan plan2 = fftw_plan_dft_r2c_1d(bothsize, freq2, freq_signal2, FFTW_ESTIMATE);
	fftw_execute(plan2);
	for (int i = 0; i < bothsize; i++)
	{
		freq1[i] = freq_signal1[i][0] * freq_signal2[i][0] + freq_signal1[i][1] * freq_signal2[i][1];
	}
	fftw_complex *freq_signal = (fftw_complex *)fftw_malloc(sizeof(fftw_complex) * bothsize);
	fftw_plan plan = fftw_plan_dft_r2c_1d(bothsize, freq1, freq_signal, FFTW_ESTIMATE);
	fftw_execute(plan);
	for (int i = 0; i < bothsize; i++)
	{
		if (freq_signal[i][0] > max_value)
		{
			max_index = i;
			max_value = freq_signal[i][0];
		}
	}
	av_free(freq1);
	av_free(freq2);
	fftw_free(freq_signal1);
	fftw_free(freq_signal2);
	fftw_free(freq_signal);
	fftw_destroy_plan(plan1);
	fftw_destroy_plan(plan2);
	fftw_destroy_plan(plan);
	return max_index * (((-1) * (*size1 < *size2)) + (*size1 >= *size2));
}
