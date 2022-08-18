/* 0ffdcae3-9db2-401b-8a82-dfad707b83cd */
/* dm01 */

$self.setName("Bronze-Arm Tribe");
$self.setPower(1000);
$self.setCivilization("nature");
$self.setRace("beast_folk");
$self.setManaCost(3);
$self.setManaReq("nature");

$self.useTrait("creature");

$match.on("card_moved", (event) => {
  if (event.cardId !== $self.id || event.to !== "battlezone") {
    return;
  }

  console.log("bronze arm tribe script test..");
});
