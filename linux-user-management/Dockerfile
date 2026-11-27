# An Amazon Linux 2023 machine for the user-management exercise.
#
# It mirrors a fresh EC2 instance: you log in as ec2-user and use sudo for
# administration. The tools the tasks need (useradd/userdel, zip/unzip, find)
# are installed explicitly because the minimal container image leaves them out.
FROM amazonlinux:2023

RUN dnf install -y shadow-utils sudo zip unzip findutils tree \
    && dnf clean all \
    && useradd -m ec2-user \
    && echo 'ec2-user ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/90-ec2-user \
    && chmod 440 /etc/sudoers.d/90-ec2-user

COPY setup.sh tasks.sh /opt/scripts/
RUN chmod +x /opt/scripts/*.sh

USER ec2-user
WORKDIR /home/ec2-user
CMD ["sleep", "infinity"]
