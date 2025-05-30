#include <libavformat/avformat.h>
#include <libavutil/file.h>

#include <ctype.h>
#include <stdio.h>
#include <string.h>

const char* get_file_extension(const char* filename)
{
	const char* dot = strrchr(filename, '.');
	if (!dot || dot == filename)
		return "";
	return dot + 1;
}

int check_file_extension(const char* filename)
{
	const char* extension = get_file_extension(filename);
	uint8_t length = strlen(extension);

	char buffer[length + 1];
	for (uint8_t i = 0; i < length; ++i)
	{
		buffer[i] = extension[i];
	}
	buffer[length] = '\0';

	const char* valid_extensions[] = { "flac", "mp2", "mp3", "opus", "aac" };
	uint8_t valid_extensions_count = sizeof(valid_extensions) / sizeof(valid_extensions[0]);

	for (uint8_t i = 0; i < valid_extensions_count; ++i)
	{
		if (strcmp(buffer, valid_extensions[i]) == 0)
		{
			return i + 1;
		}
	}

	return -1;
}

const char* get_extension_string(int extension_index)
{
	const char* valid_extensions[] = { "flac", "mp2", "mp3", "opus", "aac" };
	if (extension_index > 0 && extension_index <= sizeof(valid_extensions) / sizeof(valid_extensions[0]))
	{
		return valid_extensions[extension_index - 1];
	}
	return NULL;
}
