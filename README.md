This repo is an API to serve values based on time. These values are then used by LED Controllers to dim lights and set a specific color.

# How it works?

Brightness and color are values from 0 to 1000. This repo  uses simple linear function that increase that number from midnight to noon and then decrease it until midnight.

```
        ▲
   1000 │                                ┌─┐
        │                              ┌─┘ └─┐
        │                            ┌─┘     └─┐
        │                          ┌─┘         └─┐
        │                        ┌─┘             └─┐
        │                      ┌─┘                 └─┐
        │                    ┌─┘                     └─┐
        │                  ┌─┘                         └─┐
        │                ┌─┘                             └─┐
        │              ┌─┘                                 └─┐
        │            ┌─┘                                     └─┐
        │          ┌─┘                                         └─┐
        │        ┌─┘                                             └─┐
        │      ┌─┘                                                 └─┐
        │    ┌─┘                                                     └─┐
        │  ┌─┘                                                         └─┐
      0 │┌─┘                                                             └─┐
        └┴─────────────────────────────────────────────────────────────────┴─▶
         midnight                        noon                        midnight

```
