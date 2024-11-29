The purpose of this project is to learn Go developing it.

This app generates a map for a fictional city.
Generation is performed in several steps:

1. General shape of a city - borders;
1. Big roads;
1. Non-residential areas: parks and industrial areas;
1. Blocks and streets.

Every step is configurable.

Recomended initial values for generation:
1. Borders:
    1. Number of corners = 4 to 20
    1. Radii ratio = 1 to 1/4
    1. Additional corner point variation = 1/10 of max radius
1. Roads:
    1. Number of centers = 2 to 5
    1. Radii ≈ 1/10 of max city radius form previous paragraph, radii ratio = 1 to 1/4
    1. Road exits = 8 to 20
1. Areas are up to you
1. Blocks size > 1/20 of max radius.

Other values are not prohibited but can lead to odd results.

Usage:

go run ./ui for standalone version.

go run ./web for web server on :80 port.


Both are still under development.
