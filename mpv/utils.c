#include <mpv/client.h>
#include <stdlib.h>
#include <stdio.h>

char** makeStringArray(int length) {
	return calloc(sizeof(char*), length);
}

void setString(char** arr, int index, char* value) {
	arr[index] = value;
}

mpv_node* makeNodeList(int length) {
	return calloc(sizeof(mpv_node), length);
}

void setNodeListElement(mpv_node* values, int index, mpv_node value) {
	values[index] = value;
}