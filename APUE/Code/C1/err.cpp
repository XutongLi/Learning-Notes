/*************************************************************************
    > File Name: err.cpp
    > Author: lixutong
    > Created Time: 2020-04-19 19:06:11
    > Description: errno示例
 ************************************************************************/

#include "../include/apue.h"
#include <iostream>
#include <errno.h>
using namespace std;

int main(int argc, char *argv[])
{
    fprintf(stderr, "EACCES: %s\n", strerror(EACCES));
    errno = ENOENT;
    perror(argv[0]);
    exit(0);
}

