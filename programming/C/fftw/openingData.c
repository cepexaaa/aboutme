#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>
#include <libavutil/file.h>
#include <libavutil/opt.h>
#include <libswresample/swresample.h>

#include <stdio.h>

int openFile(char *arg, AVCodecContext **codec_context, int *sample_rate, AVFormatContext **context, int *i)
{
	AVStream *stream;
	if (avformat_open_input(context, arg, NULL, NULL))
		return 1;

	for (*i = 0; *i < (*context)->nb_streams; (*i)++)
	{
		if ((*context)->streams[*i]->codecpar->codec_type == AVMEDIA_TYPE_AUDIO)
		{
			stream = (*context)->streams[*i];
			break;
		}
	}
	int a = avformat_find_stream_info(*context, NULL);
	const AVCodec *codec = avcodec_find_decoder((*context)->streams[*i]->codecpar->codec_id);
	*codec_context = avcodec_alloc_context3(codec);
	(*codec_context)->pkt_timebase = (*context)->streams[*i]->time_base;
	avcodec_parameters_to_context(*codec_context, (*context)->streams[*i]->codecpar);
	if (avcodec_open2(*codec_context, codec, NULL) < 0)
	{
		fprintf(stderr, "Could not open codec\n");
		return 1;
	}
	*sample_rate = stream->codecpar->sample_rate;
	return 0;
}

double *process(char *arg, int *return_size, const int CHANNEL_NUMBER, const int *sample_rate, int streamIndex, AVCodecContext *codec_context, AVFormatContext *context)
{
	int64_t HzMax = *sample_rate;
	SwrContext *swr = swr_alloc();
	av_opt_set_int(swr, "in_channel_layout", (int64_t)codec_context->ch_layout.nb_channels, 0);
	av_opt_set_int(swr, "out_channel_layout", (int64_t)codec_context->ch_layout.nb_channels, 0);
	av_opt_set_int(swr, "in_sample_rate", codec_context->sample_rate, 0);
	av_opt_set_int(swr, "out_sample_rate", HzMax, 0);
	av_opt_set_sample_fmt(swr, "in_sample_fmt", codec_context->sample_fmt, 0);
	av_opt_set_sample_fmt(swr, "out_sample_fmt", AV_SAMPLE_FMT_DBLP, 0);
	swr_init(swr);
	int buffer_size = 8000;
	double *signal1 = (double *)av_malloc(buffer_size * sizeof(double));
	AVPacket packet;
	AVFrame *frame = av_frame_alloc();
	while (av_read_frame(context, &packet) >= 0)
	{
		if (packet.stream_index == streamIndex)
		{
			int response = avcodec_send_packet(codec_context, &packet);
			if (response < 0)
				break;
			while (response >= 0)
			{
				response = avcodec_receive_frame(codec_context, frame);
				if (response == AVERROR(EAGAIN) || response == AVERROR_EOF)
				{
					break;
				}
				if (response < 0)
				{
					return NULL;
				}
				double **temp = NULL;
				int est_samples = swr_get_out_samples(swr, frame->nb_samples);
				int ret = av_samples_alloc_array_and_samples(
					(uint8_t ***)&temp,
					frame->linesize,
					(int)codec_context->ch_layout.nb_channels,
					est_samples,
					AV_SAMPLE_FMT_DBLP,
					0);
				int frame_count = swr_convert(swr, (uint8_t **)temp, est_samples, (const uint8_t **)frame->data, frame->nb_samples);
				if (frame_count > 0)
				{
					if (buffer_size < (*return_size + frame_count))
					{
						buffer_size = (*return_size + frame_count) * 2;
						signal1 = (double *)av_realloc(signal1, buffer_size * sizeof(double));
					}
					memcpy(signal1 + *return_size, temp[0], frame_count * sizeof(double));
					*return_size += frame_count;
				}
			}
		}
		av_packet_unref(&packet);
	}
	if (frame->ch_layout.nb_channels < 2 && CHANNEL_NUMBER == 2)
	{
		fprintf(stderr, "In this file few channels");
		return NULL;
	}
	avcodec_free_context(&codec_context);
	av_frame_free(&frame);
	avformat_close_input(&context);
	swr_free(&swr);
	return signal1;
}
