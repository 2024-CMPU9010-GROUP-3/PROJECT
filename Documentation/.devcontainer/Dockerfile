FROM qmcgaw/latexdevcontainer:latest-full@sha256:4932d0d062665156f0b08ec963a80dcda94b17fa33c8b1fa7ae667dbf91d5c65

# Install additional packages
RUN apt update && apt install -y \
    make \
    curl \
    inkscape \
    libdrm2 \
    libgbm1 \
    libasound2

# Install d2lang
RUN curl -fsSL https://d2lang.com/install.sh | sh -s --
