create table users
(
    id UUID,
    username varchar(255) not null,
<<<<<<< HEAD
    buyOrRent varchar(64) not null,
=======
>>>>>>> ef43dad
    city varchar(128),
    rooms text[],
    typeHouse text[],

    yearBuiltFrom int,
    priceFrom int,
    floorFrom int,
    floorInTheHouseFrom int,

    areaFrom varchar(64),
    kitchenFrom varchar(64),

    yearBuiltTo int,
    priceTo int,
    floorTo int,
    floorInTheHouseTo int,
    areaTo varchar(64),
    kitchenTo varchar(64),

    notFirstFloor bool,
    notLastFloor bool,
    fromOwner bool,
    newBuilding bool,
    realEstate bool,
    running bool,

    primary key(id, username)
)