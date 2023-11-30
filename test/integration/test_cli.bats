#!/usr/bin/env bats
# Copyright (c) 2018 John Dewey

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to
# deal in the Software without restriction, including without limitation the
# rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
# sell copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
# FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
# DEALINGS IN THE SOFTWARE.

GILT_TEST_BASE_DIR=test/integration
GILT_TEST_BASE_TMP_DIR=${GILT_TEST_BASE_DIR}/tmp
GILT_TEST_CLIENT_DIR=examples/go-client
GILT_LIBRARY_DIR=${GILT_TEST_BASE_TMP_DIR}/library
GILT_ROLES_DIR=${GILT_TEST_BASE_TMP_DIR}/roles
GILT_TEST_DIR=${GILT_TEST_BASE_TMP_DIR}/tests
GILT_PROGRAM="../../../main.go"
GILT_DIR=~/.gilt/clone

setup() {
	GILT_CLONED_REPO_1=${GILT_DIR}/cache/https---github.com-retr0h-ansible-etcd.git-77a95b7
	GILT_CLONED_REPO_2=${GILT_DIR}/cache/https---github.com-retr0h-ansible-etcd.git-1.1
	GILT_CLONED_REPO_3=${GILT_DIR}/cache/https---github.com-lorin-openstack-ansible-modules.git-2677cc3

	GILT_CLONED_REPO_1_DST_DIR=${GILT_TEST_BASE_TMP_DIR}/retr0h.ansible-etcd
	GILT_CLONED_REPO_2_DST_DIR=${GILT_TEST_BASE_TMP_DIR}/retr0h.ansible-etcd-tag

  mkdir -p ${GILT_DIR}

	mkdir -p ${GILT_LIBRARY_DIR}
	mkdir -p ${GILT_ROLES_DIR}
	cp test/Giltfile.yaml ${GILT_TEST_BASE_TMP_DIR}/Giltfile.yaml
}

teardown() {
	rm -rf ${GILT_CLONED_REPO_1}
	rm -rf ${GILT_CLONED_REPO_2}
	rm -rf ${GILT_CLONED_REPO_3}
	rm -rf ${GILT_CLONED_REPO_1_DST_DIR}
	rm -rf ${GILT_CLONED_REPO_2_DST_DIR}

	rm -rf ${GILT_LIBRARY_DIR}
	rm -rf ${GILT_ROLES_DIR}
	rm -rf ${GILT_TEST_DIR}
	rm -f ${GILT_TEST_BASE_TMP_DIR}/Giltfile.yaml
}

@test "invoke gilt without arguments prints usage" {
	run go run main.go

	[ "$status" -eq 0 ]
	echo "${output}" | grep "GIT layering command line tool."
}

@test "invoke gilt version subcommand" {
	run go run main.go version

	[ "$status" -eq 0 ]
	echo "${output}" | jq '.date'
	echo "${output}" | jq '.build'
	echo "${output}" | jq '.version'
}

@test "invoke gilt overlay subcommand" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} overlay"

	[ "$status" -eq 0 ]
}

@test "invoke gilt overlay subcommand with gilt-file flag" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} --gilt-file Giltfile.yaml overlay"

	[ "$status" -eq 0 ]
}

@test "invoke gilt overlay subcommand with gilt-dir flag" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} --gilt-dir ${GILT_DIR} overlay"

	[ "$status" -eq 0 ]
}

@test "invoke gilt overlay subcommand with f flag" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} -f Giltfile.yaml overlay"

	[ "$status" -eq 0 ]
}

@test "invoke gilt overlay subcommand with c flag" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} -c ${GILT_DIR} overlay"

	[ "$status" -eq 0 ]
}

@test "invoke gilt overlay subcommand with d flag" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} -d overlay"

	[ "$status" -eq 0 ]
}

@test "invoke gilt overlay subcommand with debug flag" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} --debug overlay"

	[ "$status" -eq 0 ]
	echo "${output}" | grep "[https://github.com/retr0h/ansible-etcd.git@77a95b7]"
	echo "${output}" | grep -E ".*cloning.*https---github.com-retr0h-ansible-etcd.git-77a95b7"
}

@test "invoke gilt overlay when already cloned" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} overlay"
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} overlay"

	echo "${output}" | grep "clone already exists"
}

@test "invoke gilt overlay and clone" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} overlay"

	run stat ${GILT_CLONED_REPO_1}
	[ "$status" = 0 ]

	run stat ${GILT_CLONED_REPO_2}
	[ "$status" = 0 ]

	run stat ${GILT_CLONED_REPO_3}
	[ "$status" = 0 ]
}

@test "invoke gilt overlay and checkout index" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} overlay"

	run stat ${GILT_CLONED_REPO_1_DST_DIR}
	[ "$status" = 0 ]

	run stat ${GILT_CLONED_REPO_2_DST_DIR}
	[ "$status" = 0 ]
}

@test "invoke gilt overlay and checkout by tag" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} overlay"

	run stat ${GILT_CLONED_REPO_1_DST_DIR}
	[ "$status" = 0 ]

	run stat ${GILT_CLONED_REPO_2_DST_DIR}
	[ "$status" = 0 ]
}

@test "invoke gilt overlay and copy sources" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} overlay"

	# Copy src file matched by regexp to dst dir.
	run stat ${GILT_LIBRARY_DIR}/cinder_manage
	[ "$status" = 0 ]
	run stat ${GILT_LIBRARY_DIR}/glance_manage
	[ "$status" = 0 ]
	run stat ${GILT_LIBRARY_DIR}/heat_manage
	[ "$status" = 0 ]
	run stat ${GILT_LIBRARY_DIR}/keystone_manage
	[ "$status" = 0 ]
	run stat ${GILT_LIBRARY_DIR}/nova_manage
	[ "$status" = 0 ]

	# Copy src file to dst dir.
	run stat ${GILT_LIBRARY_DIR}/nova_quota
	[ "$status" = 0 ]

	# Copy src file to dst file.
	run stat ${GILT_LIBRARY_DIR}/neutron_router.py
	[ "$status" = 0 ]

	# Copy src dir to dst dir.
	run stat ${GILT_TEST_DIR}/keystone_service.py
	echo $output
	[ "$status" = 0 ]
	run stat ${GILT_TEST_DIR}/test_keystone_service.py
	echo $output
	[ "$status" = 0 ]
}

@test "invoke gilt overlay and post commands" {
	run bash -c "cd ${GILT_TEST_BASE_TMP_DIR}; go run ${GILT_PROGRAM} overlay"

	run stat ${GILT_TEST_BASE_TMP_DIR}/ansible-etcd-repo-post-command-1
	[ "$status" = 0 ]

	run stat ${GILT_TEST_BASE_TMP_DIR}/openstack-ansible-modules-repo-post-command-1
	[ "$status" = 0 ]

	run stat ${GILT_TEST_BASE_TMP_DIR}/openstack-ansible-modules-repo-post-command-2
	[ "$status" = 0 ]
}

@test "invoke gilt client" {
	run bash -c "cd ${GILT_TEST_CLIENT_DIR}; go mod tidy; go run main.go"
 	[ "$status" -eq 0 ]
 	echo $output

	run stat ${GILT_TEST_BASE_TMP_DIR}/retr0h.ansible-etcd/
	[ "$status" = 0 ]
}
