//
// Created by paul on 6/30/15.
//

#include <stdio.h>
#include <sys/types.h>
#include <ifaddrs.h>
#include <netinet/in.h>
#include <string.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <string>
#include <unistd.h>

#include <chrono>

#include "dautil.hpp"

long currentTimeMillis() {
    using namespace std::chrono;
    milliseconds ms = duration_cast<milliseconds>(
            system_clock::now().time_since_epoch()
    );
    return ms.count();
}

std::string getHostnameBestGuess() {
    std::string result("unknown");
    struct addrinfo hints, *info, *p;
    int gai_result;

    char hostname[1024];
    hostname[1023] = '\0';
    gethostname(hostname, 1023);

    memset(&hints, 0, sizeof hints);
    hints.ai_family = AF_UNSPEC; /*either IPV4 or IPV6*/
    hints.ai_socktype = SOCK_STREAM;
    hints.ai_flags = AI_CANONNAME;

    if ((gai_result = getaddrinfo(hostname, "http", &hints, &info)) != 0) {
        fprintf(stderr, "getaddrinfo: %s\n", gai_strerror(gai_result));
        exit(1);
    }

    //take the first one that isn't null
    for (p = info; p != NULL; p = p->ai_next) {
        if (p->ai_canonname != NULL) {
            result.assign(p->ai_canonname);
            break;
        }
        //printf("hostname: %s\n", p->ai_canonname);
    }

    freeaddrinfo(info);

    return result;
}