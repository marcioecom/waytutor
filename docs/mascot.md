# Mascote

A capivara com pilha de livros e abajur. Renderizada no Blender, exportada como
clip curto por estado. Nada de código no asset: a extensão só troca o arquivo e
aplica o loop.

## Estados

Um estado por situação real do painel. Seis clipes cobrem o produto; os quatro
primeiros são loops, os dois últimos podem ser quadro parado.

| Estado | Quando aparece | O que a capivara faz |
| --- | --- | --- |
| `idle` | sem material, ou leitura terminada | dorme na pilha, respira devagar, abajur aceso mas baixo |
| `reading` | `playing` | olhos acompanham a página, orelha balança, vira folha de vez em quando |
| `thinking` | `loading`, síntese ou resumo a caminho | olhos no abajur, cabeça imóvel, luz do abajur pulsando |
| `paused` | `paused`, conta pendente ou pausada | olha para o leitor, sem se mexer |
| `error` | erro de leitura, cota, voz | orelhas caídas, abajur fraco |
| `done` | resumo aberto | caneta na pata, escreve no caderno |

Regras de continuidade: o abajur não muda de lugar, os livros não mudam de
posição e a câmera nunca se move. Só a pose, a respiração e a luz mudam. É a
mesma mesa, sempre.

## Contrato de exportação

- Área de desenho `400 × 200` a 1×, exportada a `800 × 400` (pantalla retina),
  exibida em `400 × 200`.
- Fundo transparente. O painel aplica o fundo e a sombra; o clip não traz fundo
  nem vinheta.
- Loop de 2 a 4 segundos, 12 a 15 quadros por segundo, primeiro e último quadro
  idênticos para o loop não pular.
- Um arquivo por estado: `idle`, `reading`, `thinking`, `paused`, `error`, `done`.
- Nome do arquivo igual ao nome do estado: `reading.webp`.
- Peso máximo de 150 KB por clip. O painel inteiro hoje tem 132 KB; seis clips
  de 1 MB mudariam a cara do download na primeira instalação.
- Sem texto queimado no vídeo, sem marca d'água, sem emoji.
- A luz do abajur não pisca em frequência acima de 3 Hz.

## Integração

Os clips ficam em `apps/extension/public/mascot/<estado>.webp` e são trocados por
um componente único, que recebe o estado e mantém o clip anterior em tela até o
novo carregar. Nada de animação em CSS por cima do clip, exceto o fade de troca.
