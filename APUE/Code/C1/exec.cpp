/*************************************************************************
    > File Name: exec.cpp
    > Author: lixutong
    > Created Time: 2020-04-19 16:34:35
    > Description: 从标准输入读取命令，并执行这些命令
 ************************************************************************/

#include "apue.h"
#include <iostream>
using namespace std;

void sig_int(int signo) {
    printf("interrupt\n %% ");
}

int main()
{
    char buf[MAXLINE];
    pid_t pid;
    int status;
    
    if (signal(SIGINT, sig_int) == SIG_ERR) { //中断键处理
        cout << "signal error" << endl;
        exit(1);
    }

    printf("%% ");
    while (fgets(buf, MAXLINE, stdin) != NULL) {
        if (buf[strlen(buf) - 1] == '\n')
            buf[strlen(buf) - 1] = 0;
        if ((pid = fork()) < 0) {
            cout << "fork error" << endl;
            exit(1);
        }
        else if (pid == 0) {
            execlp(buf, buf, (char *)0);
            cout << "couldn't execute: " << buf;
            exit(127);
        }
        if ((pid = waitpid(pid, &status, 0)) < 0) {
            cout << "waitpid error" << endl;
            exit(1);
        }
        printf("%% ");
    }
    exit(0);
}
