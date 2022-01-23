#include <mpv/client.h>

mpv_node* makeNodeList(int);
void setNodeListElement(mpv_node*, int, mpv_node);

char** makeStringArray(int);
void setString(char**, int, char*);