PYTHON_VERSION:=3.9.2
PROJECT:=$(shell basename $$PWD)

venv:
	pyenv local ${PYTHON_VERSION}
	pyenv virtualenv ${PROJECT}
	pyenv local ${PROJECT}
	pip install -U pip setuptools wheel
	pip install -r requirements.txt
	ansible-galaxy install -r requirements.yaml
