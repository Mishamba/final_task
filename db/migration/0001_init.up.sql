CREATE TABLE Tweets (
    authors_ID int NOT NULL,
    tweet varchar(255) NOT NULL,
    post_time DATETIME,
);

CREATE TABLE Users (
    ID int AUTO_INCREMENT,
    username varchar(30) NOT NULL UNIQUE,
    email    varchar(30) NOT NULL UNIQUE,
    password varchar(30) NOT NULL,
);
