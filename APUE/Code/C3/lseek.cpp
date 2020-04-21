/*************************************************************************
    > File Name: lseek.cpp
    > Author: lixutong
    > Created Time: 2020-04-21 19:08:37
    > Description: 
 ************************************************************************/

#include "apue.h"
#include <iostream>
using namespace std;

int main()
{
    if (lseek(STDIN_FILENO, 0, SEEK_CUR) == -1)
        cout << "cannot seek" << endl;
    else
        cout << "seek OK" << endl;
    exit(0);
}

