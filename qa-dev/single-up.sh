#!/bin/bash
#
set -eu

# must have a paramater
if [[ -z $1 ]]; then
	echo "You should Give me a Project_Name;"
	exit 1
fi

PROJECT_NAME=$1
# 将项目名称 传入 ansible-playbook 远程创建镜像
ansible-playbook playbooks/build_images.yaml -e "PNAME=${PROJECT_NAME}"

# 同时创建容器， 并运行于supervisord中
ansible-playbook playbooks/deploy_program.yaml -e "PNAME=${PROJECT_NAME}"
