//
// Created by paul on 6/12/15.
//

#ifndef DOCKER_AGENT_DAUTIL_HPP
#define DOCKER_AGENT_DAUTIL_HPP


#include <chrono>

// ...


/**
 * derp
 */
long currentTimeMillis() {
    using namespace std::chrono;
    milliseconds ms = duration_cast<milliseconds>(
            system_clock::now().time_since_epoch()
    );
    return ms.count();
}

#endif //DOCKER_AGENT_DAUTIL_HPP
