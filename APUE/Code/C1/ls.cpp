/*************************************************************************
    > File Name: ls.cpp
    > Author: lixutong
    > Created Time: 2020-04-19 05:51:04
    > Description: C++ implement of ls
 ************************************************************************/

#include "../include/apue.h"
#include <iostream>
#include <dirent.h>
using namespace std;

int main(int argc, char *argv[]) 
{
    DIR *dp;
    struct dirent *dirp;

    if (argc != 2)
        cout << "usage: ls direction_name" << endl;

    if ((dp = opendir(argv[1])) == nullptr)
        cout << "Can't open %s" << argv[1] << endl;
    while ((dirp = readdir(dp)) != nullptr)
        cout << dirp->d_name << endl;

    closedir(dp);
    exit(0);
}

