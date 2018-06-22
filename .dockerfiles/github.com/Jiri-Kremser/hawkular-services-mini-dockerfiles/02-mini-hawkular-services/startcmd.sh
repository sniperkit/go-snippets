#!/bin/sh
#
# Copyright 2016 Red Hat, Inc. and/or its affiliates
# and other contributors as indicated by the @author tags.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

echo ${HOSTNAME} > /etc/machine-id
${AS_ROOT}/bin/add-user.sh -a -u ${HAWKULAR_USER} -p ${HAWKULAR_PASSWORD} -g read-write,read-only
${AS_ROOT}/bin/standalone.sh -b 0.0.0.0 \
       -bmanagement 0.0.0.0 \
       -Djboss.server.name=${HOSTNAME} \
       -Djboss.server.data.dir=/opt/data/data \
       -Djboss.server.log.dir=/opt/data/log \
       -Dactivemq.artemis.client.global.thread.pool.max.size=${HAWKULAR_JMS_THREAD_POOL:-30} \
       -Dhawkular.agent.enabled=${HAWKULAR_AGENT_ENABLE} \
       -Dhawkular.rest.user=${HAWKULAR_USER} \
       -Dhawkular.rest.password=${HAWKULAR_PASSWORD} \
       -Dhawkular.metrics.default-ttl=${HAWKULAR_METRICS_TTL:-14}
