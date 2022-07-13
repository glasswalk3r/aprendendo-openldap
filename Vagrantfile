# frozen_string_literal: true

Vagrant.configure('2') do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  # config.vm.box = "roboxes/centos7"
  config.vm.box = 'ARFREITAS/centos7'

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  config.vm.box_check_update = false

  config.vm.provider 'virtualbox' do |vb|
    vb.memory = 1024
    vb.name = 'ldap-base'
    vb.linked_clone = true
  end

  config.vm.provision 'shell', inline: 'yum makecache fast && yum upgrade -y && yum install python3 -y'

  config.vm.provision :ansible do |ansible|
    ansible.playbook = 'ntp.yaml'
    ansible.config_file = 'ansible.cfg'
  end

  config.vm.define 'vinfra' do |i|
    i.vm.provider 'virtualbox' do |vb|
      vb.name = 'vinfra'
    end
    i.vm.hostname = 'vinfra.local.br'
    i.vm.network 'private_network', ip: '192.168.56.85'

    i.vm.provision :ansible do |ansible|
      ansible.playbook = 'vinfra/main.yaml'
      ansible.config_file = 'ansible.cfg'
    end
  end

  config.vm.define 'master' do |m|
    m.vm.provider 'virtualbox' do |vb|
      vb.name = 'master'
    end
    m.vm.hostname = 'master.local.br'
    m.vm.network 'private_network', ip: '192.168.56.80'

    m.vm.provision :ansible do |ansible|
      ansible.playbook = 'master/main.yaml'
      ansible.config_file = 'ansible.cfg'
    end
  end
end

# -*- mode: ruby -*-
# vi: set ft=ruby :
