CREATE TABLE category (
    id          INTEGER NOT NULL,
    name        VARCHAR,
    description VARCHAR NOT NULL,
    PRIMARY KEY (
        id
    ),
    UNIQUE (
        name
    )
);

CREATE TABLE question (
    id          INTEGER NOT NULL,
    ticket      INTEGER NOT NULL,
    number      INTEGER NOT NULL,
    title       VARCHAR NOT NULL,
    help        VARCHAR NOT NULL,
    image       VARCHAR NOT NULL,
    category_id VARCHAR,
    PRIMARY KEY (
        id
    ),
    FOREIGN KEY (
        category_id
    )
    REFERENCES category (name)
);

CREATE TABLE answer (
    id             INTEGER NOT NULL,
    title          VARCHAR NOT NULL,
    correct_answer BOOLEAN NOT NULL,
    question_id    INTEGER,
    PRIMARY KEY (
        id
    ),
    FOREIGN KEY (
        question_id
    )
    REFERENCES question (id)
);
