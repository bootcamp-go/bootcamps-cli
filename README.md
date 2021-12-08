# bootcamps

## Install

To install the CLI just:

```sh
go get github.com/ezedh/bootcamps
```

## Usage

### Configure CLI

To configure the CLI you will require these variables:

- token: Personal Github Api Token.
- username: Github username.
- company: Company of the current bootcamp.

#### Basic usage

```sh
bootcamps configure
```

This will ask you for each variable in order.

#### Set token

```sh
bootcamps configure token
```

#### Set DH username

```sh
bootcamps configure username
```

#### Set DH company

```sh
bootcamps configure company
```

## Create Repos

In order to create repos you will first need to:

- Have a template in the branch of the configured company inside templates repo.
- Have a wave config file in the branch of the configured company inside users repo.

### Usage

```sh
bootcamps create
```

The CLI will ask you the wave. This wave **must** have its config file.

## Config file

The config file of the wave consist of a .yaml file inside users repo (private) with these parameters:

- teachers: Objet where each key is the teacher's github username and the value is an array of the groups number assigned.
- groups: Object where each key is the group number and the value is an array of the group members.

**group numbers must be strings**

### Example

```yaml
teachers:
  ezegrosfeld:
    - '1'
    - '4'
  ezedh:
    - '2'
    - '3'
groups:
  '1':
    - idontexist
    - octokit
    - testinguser
  '2':
    - octokit
    - testinguser
  '3':
    - idontexist
    - octokit
  '4':
    - idontexist
    - testinguser
```

The configuration above (for 4 groups) will create 4 repositories.

## Get Sprint Repos

To get a sprint repo you will use the wave, the group, the sprint and the owner of the repo using flags:

```sh
bootcamps get -w 10 -g 1 -s 1 -o ezedh
```

Also works:

```sh
bootcamps get --wave 10 --group 1 --sprint 1 -o ezedh
```

## Test Sprint

To test a sprint you will use the wave, the group, the sprint and the owner of the repo using flags:

```sh
bootcamps test -w 10 -g 1 -s 1 -o ezedh
```

Also works:

```sh
bootcamps test --wave 10 --group 1 --sprint 1 --owner ezedh
```
