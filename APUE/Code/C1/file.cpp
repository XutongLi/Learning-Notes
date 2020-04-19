/*************************************************************************
    > File Name: file.cpp
    > Author: lixutong
    > Created Time: 2020-04-20 04:13:03
    > Description: 
 ************************************************************************/

#include "../include/apue.h"
#include <iostream>
using namespace std;

int main()
{   
    printf("uid = %d, gid = %d\n", getuid(), getgid());
    exit(0);
}

