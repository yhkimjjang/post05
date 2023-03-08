drop table if exists Users;
drop table if exists Userdata;

create table Users (
	ID serial,
	Username varchar(100) primary key
);

create table Userdata (
	UserID int not null,
	Name varchar(100),
	Surname varchar(100),
	Descrption varchar(200)
);