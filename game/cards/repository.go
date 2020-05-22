package cards

import (
	"duel-masters/game/cards/dm01"
	"duel-masters/game/match"
)

// Cards is a map with all the card id's in the game and corresponding CardConstructor
var Cards = map[string]match.CardConstructor{

	// dm01
	"57eeb3c3-2561-4841-a381-2e50d17533d1": dm01.AquaHulcus,
	"ecd1ae69-4f63-4e8d-a3f4-9a5c81f98a20": dm01.EmeraldGrass,
	"09b218fc-9c5a-48ef-9555-4908932271e9": dm01.AquaKnight,
	"4097a036-a775-4218-9a1d-f57ead85dda6": dm01.AquaSniper,
	"c43bc627-9e7a-4686-9d61-789425669b02": dm01.AquaSoldier,
	"9781089f-1aa9-4a75-b106-35e9d431e31d": dm01.AquaVehicle,
	"1d72eb3e-5185-449a-a16f-391bd2338343": dm01.BurningMane,
	"fcd0cb50-b687-4180-90a8-390aeb8705cc": dm01.FearFang,
	"10e0e90f-ad7d-4b69-98d5-f01525eb1cdd": dm01.SteelSmasher,
	"015fd6bb-37a9-45cf-bb6b-a5497412b880": dm01.BronzeArmTribe,
	"6663848d-035e-44b6-9d9f-7b236ea5bc43": dm01.GoldenWingStriker,
	"0e26fe1a-a9d1-4c78-80e9-7f4cc0e4c1c8": dm01.MightyShouter,
	"0b1e4f56-6342-46db-9faf-882fd1f1f179": dm01.ArtisanPicora,
	"983e72d7-3f4e-466d-a4e3-06552e392af2": dm01.NomadHeroGigio,
	"0cc5279e-0a26-41a8-a2a5-f7711120b772": dm01.LahPurificationEnforcer,
	"808ddd60-e8ca-49f0-9baa-57e632f85b28": dm01.RaylaTruthEnforcer,
	"91db2302-6794-4aa4-b17b-6637d356e9ac": dm01.AstrocometDragon,
	"0ffdcae3-9db2-401b-8a82-dfad707b83cd": dm01.BolshackDragon,
	"6cf85053-abaa-4577-b151-86123004980e": dm01.Draglide,
	"3b6e6c29-017d-41b9-bf93-186f7963723e": dm01.GatlingSkyterror,
	"1c5511be-7629-41c5-bf17-4bc810be5472": dm01.ScarletSkyterror,
	"a4adb373-0aec-4fff-997c-3820c7ec528d": dm01.DomeShell,
	"1ecb54a2-bcbf-4396-bf09-50dfe984e287": dm01.StormShell,
	"c761c174-87c3-4f4a-ab94-aa837c5ab587": dm01.TowerShell,
	"18e0e199-7827-4a4c-a37d-3acfa4e500d6": dm01.RoaringGreatHorn,
	"2aeae452-5630-4f86-b073-7e9dc07adc43": dm01.StampedingLonghorn,
	"84e1b416-c2d5-4ae1-aca0-025651c6aa58": dm01.TriHornShepherd,
	"3e2940f4-5654-4456-bfc2-fa5e43911cfb": dm01.KingCoral,
	"cd13f7c2-aa5e-43b8-8811-700f230a5de5": dm01.KingDepthcon,
	"f04feb7f-971f-4192-893a-46c23180233a": dm01.KingRippedHide,
	"596f5b72-2502-4120-81f9-9ff9a17271d8": dm01.CandyDrop,
	"a3cf18f0-b04f-45e9-97f7-2a2ead0a1787": dm01.FaerieChild,
	"3f331274-f5f8-42e7-9f28-ce637add34d4": dm01.MarineFlower,
	"ce48ff2c-ea9e-4c12-8629-028d2480b063": dm01.IllusionaryMerfolk,
	"4b021e6f-39cf-401e-89cf-f164f7c0a797": dm01.PhantomFish,
	"cfe9f5b8-2eeb-42c9-89ff-7e69734adc4d": dm01.RevolverFish,
	"70e6cc2c-c63d-4dd9-9b6e-0713fed174bb": dm01.SaucerHeadShark,
	"4c9acf76-cc52-44c3-9e39-613d744c63c5": dm01.PoisonousMushroom,
	"cc9762c3-515a-4734-a3fe-1e0c4c3b3d71": dm01.BoneAssassin,
	"4d3201e8-0d9b-481e-b8e3-86cb90058e20": dm01.BoneSpider,
	"ec46daa1-49ce-4b88-b2bc-e923672ad0f3": dm01.SkeletonSoldierTheDefiled,
	"90b2ed59-828c-4237-ac2e-b7008a02ad2e": dm01.WanderingBraineater,
	"5d3d7052-e5fa-4502-8d31-c72673232317": dm01.HanusaRadianceElemental,
	"6a4270cf-f3be-4c66-8b30-eb2c769065dc": dm01.IocantTheOracle,
	"25a2af16-cc42-4f4c-8c3d-59fb3a7ca74b": dm01.UrthPurifyingElemental,
	"c4839847-e393-47b0-b172-95531aa6d39e": dm01.Gigaberos,
	"5d73062e-acff-47e6-b49a-c0bb1a1762b5": dm01.Gigagiele,
	"6161e271-5294-4073-94d2-b9c06f9d8fa3": dm01.Gigargon,
	"dc1b51b3-52e7-4f1c-8770-515d4e1cb53d": dm01.DeathligerLionOfChaos,

	// dm02
	// ...
}
