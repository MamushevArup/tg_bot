create table users
(
    id UUID,
    username varchar(255) not null,
    buyOrRent varchar(64) not null,
    typeItem varchar(64) not null,
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

    primary key(id, username)
)