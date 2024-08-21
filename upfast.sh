#!/bin/bash

source .env

get_public_dns() {
    if [ -f "hosts.ini" ]; then
        awk '/\[tf2_server\]/{getline; print $1}' hosts.ini
    else
        echo "Error: hosts.ini not found. Has the server been created?" >&2
        exit 1
    fi
}

create_server() {
    if ! terraform apply; then
        echo "Error: Terraform apply failed or was declined."
        exit 1
    fi

    # Extract output values
    public_ip=$(terraform output -raw instance_public_ip)
    public_dns=$(terraform output -raw instance_public_dns)

    if ! printf '%s\n' \
        "[tf2_server]" \
        "$public_dns ansible_user=ec2-user ansible_ssh_private_key_file=$SSH_PRIVATE_KEY_PATH" \
        "" \
        "[tf2_server:vars]" \
        "srcds_token=$SRCDS_TOKEN" \
        "srcds_rconpw=$SRCDS_RCONPW" \
        "srcds_pw=$SRCDS_PW" \
        "srcds_startmap=$SRCDS_STARTMAP" \
        "srcds_hostname=$SRCDS_HOSTNAME" \
        "srcds_secured=$SRCDS_SECURED" \
        > hosts.ini
    then
        echo "Error: Failed to create hosts.ini file."
        exit 1
    fi

    if ! ansible -i hosts.ini -m ping $public_dns; then
        echo "Error: Ansible ping failed. Check server reachability and configuration."
        exit 1
    fi

    echo "Server is reachable and configured correctly."
    echo "IP: $public_ip"
}

start_server() {
    ansible-playbook -i hosts.ini tf2_server_playbook.yml
}

restart_server() {
    ansible-playbook -i hosts.ini tf2_server_playbook.yml --tags restart
}

destroy_server() {
    if terraform destroy; then
        rm hosts.ini
        echo "Server destroyed and hosts.ini cleared."
    else
        echo "Error: Failed to destroy server."
        exit 1
    fi
}

connect_to_server() {
    public_dns=$(get_public_dns)
    ssh -i "$SSH_PRIVATE_KEY_PATH" ec2-user@"$public_dns"
}

connect_to_container() {
    public_dns=$(get_public_dns)
    ssh -i "$SSH_PRIVATE_KEY_PATH" ec2-user@"$public_dns" -t "sudo docker attach tf2-dedicated"
}

usage() {
    echo "Usage: $0 [create|start|restart|destroy|connect|connect_container]"
    echo "  create  - Create a new server"
    echo "  start   - Start an existing server"
    echo "  restart - Restart the server"
    echo "  destroy - Destroy the server"
    echo "  connect - Connect to the server"
    echo "  connect_container - Connect to the container"
}

case "$1" in
    create)
        create_server
        ;;
    start)
        start_server
        ;;
    restart)
        restart_server
        ;;
    destroy)
        destroy_server
        ;;
    connect)
        connect_to_server
        ;;
    connect_container)
        connect_to_container
        ;;
    *)
        usage
        exit 1
        ;;
esac