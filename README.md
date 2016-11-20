# ChatINSPER -- Tecnologias Web

### Pedro Cunial

Terceiro e final projeto da diciplina de tecnologias web do curso de engenharia de computação no INSPER

## Config

- Tenha suas variáveis de sistema referentes ao Go, assim como PATH no seu profile (.bashrc, .zshrc etc), por exemplo:

```
export GOPATH=$HOME/Documents/Go
export PATH="$PATH:$GOPATH/bin"

export PORT=8888
```

- Instale o Godep, ele cuidará das dependencias do projeto por você. Você pode instalar o Godep utilizando o comando:

```
go get gihub.com/tools/godep
```

- Monte o projeto!

```
go build
```

- Um arquivo executavel deve ter aparecido no diretório, execute ele e acesse localhost:<PORT>, onde PORT será o valor padrão da sua variável de sistema $PORT

- Navegue para web/views/chat.html e utilize o serviço

O serviço pode ser acessado também pelo site:
https://shrouded-springs-35758.herokuapp.com/view/chat.html


