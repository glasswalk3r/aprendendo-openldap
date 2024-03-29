# frozen_string_literal: true

dns_domain = 'local.br'
base_dn = 'dc=local,dc=br'
master_ldap_server = "master.#{dns_domain}"
slave_ldap_server = "slave.#{dns_domain}"
admin_pass = '123456'
sync_password = '654321'
admin_dn = "cn=Manager,#{base_dn}"
replication_user = 'rpuser'

Vagrant.configure('2') do |config|
  config.vm.box = 'roboxes/centos7'

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  config.vm.box_check_update = true

  config.vm.provider 'virtualbox' do |vb|
    vb.memory = 1024
    vb.name = 'ldap-base'
    vb.linked_clone = true
    # to avoid warnings from Virtualbox
    vb.customize ['modifyvm', :id, '--vrde', 'off']
    vb.customize ['modifyvm', :id, '--graphicscontroller', 'vmsvga']
  end

  config.vm.provision 'shell', inline: 'yum makecache fast && yum upgrade -y && yum install python3 tree mailx -y'

  config.vm.provision :ansible do |ansible|
    ansible.playbook = 'network.yaml'
    ansible.config_file = 'ansible.cfg'
  end

  config.vm.define 'vinfra' do |i|
    i.vm.provider 'virtualbox' do |vb|
      vb.name = 'vinfra'
    end
    i.vm.hostname = "vinfra.#{dns_domain}"
    i.vm.network 'private_network', ip: '192.168.56.85'

    i.vm.provision :ansible do |ansible|
      ansible.playbook = 'vinfra/main.yaml'
      ansible.config_file = 'ansible.cfg'
    end
  end

  # enable debugging Ansible configuring tasks
  # ENV['ANSIBLE_VERBOSITY'] = '3'

  # DRY
  openldap_role_vars = {
    'admin_pass' => admin_pass,
    'admin_dn' => admin_dn,
    'base_dn' => base_dn
  }

  config.vm.provision :ansible do |ansible|
    ansible.playbook = 'openldap.yaml'
    ansible.config_file = 'ansible.cfg'
    ansible.host_vars = {
      'master' => openldap_role_vars,
      'slave' => openldap_role_vars
    }
  end

  # DRY
  master_slave_vars = {
    'admin_pass' => admin_pass,
    'sync_password' => sync_password,
    'admin_dn' => admin_dn,
    'base_dn' => base_dn,
    'replication_user' => replication_user
  }

  config.vm.define 'master' do |m|
    m.vm.provider 'virtualbox' do |vb|
      vb.name = 'master'
    end
    m.vm.hostname = master_ldap_server
    m.vm.network 'private_network', ip: '192.168.56.80'
    m.vm.network 'forwarded_port', guest: 389, host: 3389, host_ip: '127.0.0.1', id: 'ldap'

    m.vm.provision :ansible do |ansible|
      ansible.playbook = 'master/main.yaml'
      ansible.config_file = 'ansible.cfg'
      ansible.host_vars = {
        'master' => master_slave_vars
      }
    end
  end

  config.vm.define 'slave' do |m|
    m.vm.provider 'virtualbox' do |vb|
      vb.name = 'slave'
    end
    m.vm.hostname = slave_ldap_server
    m.vm.network 'private_network', ip: '192.168.56.81'
    m.vm.network 'forwarded_port', guest: 389, host: 4389, host_ip: '127.0.0.1', id: 'ldap'

    # m.vm.provision :ansible do |ansible|
    #   ansible.playbook = 'slave/main.yaml'
    #   ansible.config_file = 'ansible.cfg'
    #   ansible.host_vars = {
    #     'master' => master_slave_vars
    #   }
    # end
  end

  config.vm.define 'client' do |c|
    c.vm.provider 'virtualbox' do |vb|
      vb.name = 'client'
    end
    c.vm.hostname = "client.#{dns_domain}"
    c.vm.network 'private_network', ip: '192.168.56.82'

    c.vm.provision :ansible do |ansible|
      ansible.playbook = 'client/main.yaml'
      ansible.config_file = 'ansible.cfg'
      ansible.host_vars = {
        'client' => {
          'base_dn' => base_dn,
          'master_ldap_server' => master_ldap_server
        }
      }
    end
  end
end

# -*- mode: ruby -*-
# vi: set ft=ruby :
