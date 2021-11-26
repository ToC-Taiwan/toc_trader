# CHANGELOG

## v1.4.0 (2021-11-26)

### New feature

- **simulate**: :truck: change flow to debug can choose n, add round to rsi and ticks period limit([`3acb824`](https://gitlab.tocraw.com/root/toc_trader/commit/3acb824d0c9d29f32a5d8dd3c0a8684b17a6cbe5)) (@TimHsu@M1BP-20210907)
- **balance**: :beers: add balance, resualt, cond api, all gorm.Model add json tag, modify some api([`eca7f40`](https://gitlab.tocraw.com/root/toc_trader/commit/eca7f40916753359834137b31f552c98bd3da3fa)) (@TimHsu@M1BP-20210907)
- **simulate**: :pushpin: change the lowest tick per second to 20([`c9f0357`](https://gitlab.tocraw.com/root/toc_trader/commit/c9f0357875d4018063a15d570f19fc7554579839)) (@TimHsu@M1BP-20210907)

### Bugs fixed

- **logger**: :building_construction: fix wrong log file name([`3ff9ac1`](https://gitlab.tocraw.com/root/toc_trader/commit/3ff9ac100bd39897fde69339b79ce5008409ca4c)) (@TimHsu@M1BP-20210907)

## v1.3.0 (2021-11-24)

### New feature

- **entiretick**: :beers: add sinopac native tick_type in entiretick and tse_entiretick([`dc199bd`](https://gitlab.tocraw.com/root/toc_trader/commit/dc199bd9983a3c0bf38da9ab2bc8969e13883712)) (@TimHsu@M1BP-20210907)
- **tradebot**: :boom: clear close array if need sell or buy later, result order by balance/trade_count([`38604c9`](https://gitlab.tocraw.com/root/toc_trader/commit/38604c9d05ba16826c3a1659bc75a785d2e6143f)) (@TimHsu@M1BP-20210907)
- **stockutil**: :card_file_box: add get max, min by open used in get price, modify ci logs path([`ea727d3`](https://gitlab.tocraw.com/root/toc_trader/commit/ea727d3ec5b0950ee26cf87af64bf49df9d1d3e6)) (@TimHsu@M1BP-20210907)
- **tradebot**: :tada: remove all recover, extend normal wait time to 5 minutes, remove rsi < 50 in sell first check([`734b2bb`](https://gitlab.tocraw.com/root/toc_trader/commit/734b2bbe53d856021e308f59fc24aca8686b8c21)) (@TimHsu@M1BP-20210907)
- **simulate**: :seedling: add use default, modify get sell buy later method, remove max hold time([`e30d7f5`](https://gitlab.tocraw.com/root/toc_trader/commit/e30d7f59849b13200d4e7bcc3e3736e1dd512631)) (@TimHsu@M1BP-20210907)
- **balance**: :bulb: add balance schema, disable default Transaction,  modify forward, reverse method([`1446d7e`](https://gitlab.tocraw.com/root/toc_trader/commit/1446d7eccfb6111c9e595e1f83a271a79a10868c)) (@TimHsu@M1BP-20210907)
- **simulate**: :pencil2: add method to delete all not best result and cond([`1b2fa7f`](https://gitlab.tocraw.com/root/toc_trader/commit/1b2fa7f74c1b15b8a0602d9d842ec45b0f45f23e)) (@TimHsu@M1BP-20210907)
- **simulate**: :package: add trade day in result, and get by trade day, modify sell, buy later time unit to 15 minutes([`4a5eaaf`](https://gitlab.tocraw.com/root/toc_trader/commit/4a5eaaf021cdc4cdb4e5d87762eed8f4093ed52a)) (@TimHsu@M1BP-20210907)
- **cron**: :pushpin: restart_sinopac_toc_trader_cron add 2, 16([`880c1a1`](https://gitlab.tocraw.com/root/toc_trader/commit/880c1a13cd817f488547a4c9ee25e1fb754864b6)) (@TimHsu@M1BP-20210907)
- **tradebot**: :heavy_plus_sign: shorten sell buy later time unit to 10 minutes([`a0da15a`](https://gitlab.tocraw.com/root/toc_trader/commit/a0da15ae088c10f8d9d28f7519f5fed90c0cd757)) (@TimHsu@M1BP-20210907)
- **simulate**: :boom: add auto simulate on startup, if best cond exist, then skip([`d1b8d2b`](https://gitlab.tocraw.com/root/toc_trader/commit/d1b8d2bd741fb3b7170e054746903229dc494f62)) (@TimHsu@M1BP-20210907)
- **tradebot**: :zap: add max hold time to rapid trade, tmp fix on 1([`c6c2e24`](https://gitlab.tocraw.com/root/toc_trader/commit/c6c2e2425d0cf9dd4722e77406cc4b80eba88e84)) (@TimHsu@M1BP-20210907)
- **tradebot**: :busts_in_silhouette: change reverse flow, add analyzeTick.Rsi < 50([`6ffdb21`](https://gitlab.tocraw.com/root/toc_trader/commit/6ffdb21f47790ccf48acc5c1f14437f269c6a4e5)) (@TimHsu@M1BP-20210907)
- **tradebot**: :zap: add stop point at 1.5% on buy, 1 % on sell first, tmp disable sell first([`80ac882`](https://gitlab.tocraw.com/root/toc_trader/commit/80ac88211adcc2724dfed41c636b601c83c48c71)) (@TimHsu@M1BP-20210907)
- **simulation**: :truck: add total loss in result, cancel ask update basic([`68a4f80`](https://gitlab.tocraw.com/root/toc_trader/commit/68a4f80a42ee7d59efe85472309e9c97d8471336)) (@TimHsu@M1BP-20210907)
- **tickanalyze**: :heavy_plus_sign: change part from 10 to 11, clear simulate condition log([`a05092e`](https://gitlab.tocraw.com/root/toc_trader/commit/a05092e8af4f06372c3c4d54a8a291332fb2073b)) (@TimHsu@M1BP-20210907)
- **simtrade**: :loud_sound: add simtrade collector to reduce log, logger add quote, add bestresult([`81d4ba7`](https://gitlab.tocraw.com/root/toc_trader/commit/81d4ba7fe0a975f9e55cff225fe2e5dd7b8643b2)) (@TimHsu@M1BP-20210907)
- **missingtick**: :iphone: finish tradeswitch for missing ticks, remove all streamticks in start([`f4796cd`](https://gitlab.tocraw.com/root/toc_trader/commit/f4796cd61c7d892aa3388aa08c39c78f82ab285e)) (@TimHsu@M1BP-20210907)
- **tradebot**: :page_facing_up: fill missing ticks when subscribe, rename some package([`48868ed`](https://gitlab.tocraw.com/root/toc_trader/commit/48868edff2febc20047a22b440856d3cfd3cf2e0)) (@TimHsu@M1BP-20210907)
- **middlebot**: :children_crossing: let stock can be buy and sell first at one day([`10abb48`](https://gitlab.tocraw.com/root/toc_trader/commit/10abb48dce06a300d3303bf151ed55b08060ec37)) (@TimHsu@M1BP-20210907)
- **tradebot**: :monocle_face: InitStartUpQuota separate buy and sell first, save unfinished stock to current map([`1ed17c3`](https://gitlab.tocraw.com/root/toc_trader/commit/1ed17c38df060492c51592e7a6cbff7c22b69c60)) (@TimHsu@M1BP-20210907)
- **tradebot**: :speech_balloon: stock can be buy and sell first at one day, modifiy cond, simulation([`75dae10`](https://gitlab.tocraw.com/root/toc_trader/commit/75dae1048de0f9d5a5bc93cb0e10a2b6d4d47642)) (@TimHsu@M1BP-20210907)

### Bugs fixed

- **logger**: :construction: change log file name time to RFC339([`1f1792b`](https://gitlab.tocraw.com/root/toc_trader/commit/1f1792be7e610213355528726d01e3156a7206e5)) (@TimHsu@M1BP-20210907)
- **ci**: :art: remove folder structure in tar([`121cb5c`](https://gitlab.tocraw.com/root/toc_trader/commit/121cb5c4f0638a4fb74d6a98cf9e295775d32f34)) (@TimHsu@M1BP-20210907)
- **ci**: :speech_balloon: change deployer to root, modify logs path([`8a7c87e`](https://gitlab.tocraw.com/root/toc_trader/commit/8a7c87e00278fd22912d16f92d8d7c9e0e85b3fa)) (@TimHsu@M1BP-20210907)
- **ci**: :mute: change tar command arguments, modify file path([`c0ab447`](https://gitlab.tocraw.com/root/toc_trader/commit/c0ab4471b9d2b77e335a8cb7897230dc545ba58b)) (@TimHsu@M1BP-20210907)
- **ci**: :wrench: fix wrong logs compressed([`5ec9210`](https://gitlab.tocraw.com/root/toc_trader/commit/5ec9210451eec898fde38a0ad4744ba654010ea6)) (@TimHsu@M1BP-20210907)
- **lastclose**: :passport_control: fix no last close panic, change panic to logrus panic, add necessary recover([`b0eec6c`](https://gitlab.tocraw.com/root/toc_trader/commit/b0eec6cf4f6518c99f4135fc16470eb465e2728a)) (@TimHsu@M1BP-20210907)
- **tradebothandler**: :truck: check if channel exists, remove some non-use method, add task recover([`63349c2`](https://gitlab.tocraw.com/root/toc_trader/commit/63349c2fc89f00b8db3f10f1caac4c46a9dbcd04)) (@TimHsu@M1BP-20210907)
- **simulate**: :tada: modify tick period count from 4 to 2, volume per second from 90 to 30([`1136067`](https://gitlab.tocraw.com/root/toc_trader/commit/1136067000b99591731e1967870f64d627fb4e66)) (@TimHsu@M1BP-20210907)
- **simulate**: :package: remove use global, modify time limit method, modify ticks per second to 100-20([`7a993a8`](https://gitlab.tocraw.com/root/toc_trader/commit/7a993a88372f349dedee1aa02d7171032ae926ed)) (@TimHsu@M1BP-20210907)
- **ip**: :monocle_face: send trader ip every startup([`2bb4f8e`](https://gitlab.tocraw.com/root/toc_trader/commit/2bb4f8e8a9324d6fe6184512ab7d18559cb3de22)) (@TimHsu@M1BP-20210907)
- **healthcheck**: :triangular_flag_on_post: remove full_restart api, rename package, variable names([`8cc39b9`](https://gitlab.tocraw.com/root/toc_trader/commit/8cc39b9a48155713d5915451bb6004b3664c8efc)) (@TimHsu@M1BP-20210907)
- **simulate**: :wheelchair: fix wrong trade day in get best simulate result([`6cc1fbf`](https://gitlab.tocraw.com/root/toc_trader/commit/6cc1fbf4ceb7f00477539d6b75dc8943126e47a8)) (@TimHsu@M1BP-20210907)
- **tradebot**: :globe_with_meridians: add trade condition query api, separate get simulate best cond([`76ca145`](https://gitlab.tocraw.com/root/toc_trader/commit/76ca1458688e066853a8d703db16dbd5884901d1)) (@TimHsu@M1BP-20210907)
- **api_handler**: :sparkles: fix some tiny bug, hide progress bar in docker([`3ae75f1`](https://gitlab.tocraw.com/root/toc_trader/commit/3ae75f1d8b42da43fd9e91e841f018ef5e6d1bf2)) (@TimHsu@M1BP-20210907)
- **tradebot**: :camera_flash: change max hold time unit to 20 minute, remove spare cond([`f041809`](https://gitlab.tocraw.com/root/toc_trader/commit/f041809fdeeaf1745753f58d9267dbf05f0faa73)) (@TimHsu@M1BP-20210907)
- **tradeswitch**: :arrow_up: fix sell first wrong setting([`aa8a9e9`](https://gitlab.tocraw.com/root/toc_trader/commit/aa8a9e9554abad2f67e03b99fda23aec3f3d2245)) (@TimHsu@M1BP-20210907)
- **ticks**: :twisted_rightwards_arrows: fix missing ticks switch off in reverse([`8d2df95`](https://gitlab.tocraw.com/root/toc_trader/commit/8d2df9582b629424df41751bacde82c769162f92)) (@TimHsu@M1BP-20210907)
- **ci**: :truck: fix ci build error([`22662a5`](https://gitlab.tocraw.com/root/toc_trader/commit/22662a5bfd5b97f389abd674514287ea0cad3b03)) (@TimHsu@M1BP-20210907)
- **quota**: :recycle: fix 100% cpu rate, change quota to buy and sell first([`033e349`](https://gitlab.tocraw.com/root/toc_trader/commit/033e3499c47dec2d453eb097921ff728a85baac1)) (@TimHsu@M1BP-20210907)
- **tradebot**: :globe_with_meridians: fix wrong check sell or buylater map again([`59a6274`](https://gitlab.tocraw.com/root/toc_trader/commit/59a6274cd3c366dbb50d82c8f00a012b3c367bab)) (@TimHsu@M1BP-20210907)
- **tradebot**: :boom: fix wrong check sell or buylater map([`c07b7b2`](https://gitlab.tocraw.com/root/toc_trader/commit/c07b7b21d33737ad6ca7746613b99b90046781ba)) (@TimHsu@M1BP-20210907)
- **tradebot**: :globe_with_meridians: recover one day buy sell first one time([`aa2b192`](https://gitlab.tocraw.com/root/toc_trader/commit/aa2b192f3209d8b5fcddb6939577de4e320761a8)) (@TimHsu@M1BP-20210907)
- **cond**: :bento: fix wrong cond([`900b6f3`](https://gitlab.tocraw.com/root/toc_trader/commit/900b6f3a1b7389653598b70204d74734842082ff)) (@TimHsu@M1BP-20210907)

## v1.2.0 (2021-11-02)

### New feature

- **tradebot**: :necktie: GetRSIStatus to decide buy sell([`fdc7e58`](https://gitlab.tocraw.com/root/toc_trader/commit/fdc7e58e5dd0e8bfb400913b60a3c2d933ac0671)) (@TimHsu@M1BP-20210907)
- **tradebot**: :fire: remove closediff, change reverse rsi method, change total volume limit([`21293cb`](https://gitlab.tocraw.com/root/toc_trader/commit/21293cb1de1bd7a1268643ad20e343354d14bff3)) (@TimHsu@M1BP-20210907)
- **healthcheck**: :beers: add server token check([`eb7e808`](https://gitlab.tocraw.com/root/toc_trader/commit/eb7e8085dfdbc78361c5130cc1a38e585d37c051)) (@TimHsu@M1BP-20210907)
- **tradebot**: :fire: separate forward, reverse cond, fix status, add trim historyClose switch([`83bd6ec`](https://gitlab.tocraw.com/root/toc_trader/commit/83bd6ec6562e58bec03706ab0f52d6bc9ee59523)) (@TimHsu@M1BP-20210907)
- **tradebot**: :fire: count positive days, shrink trade in time([`7bb62e6`](https://gitlab.tocraw.com/root/toc_trader/commit/7bb62e611de002f295e9d8219b331ce1036846a8)) (@TimHsu@M1BP-20210907)
- **md**: :card_file_box: add changelog and contributing([`2409dcf`](https://gitlab.tocraw.com/root/toc_trader/commit/2409dcfda00e2d1851f85f9070afaca76cb85f18)) (@TimHsu@M1BP-20210907)
- **tradebot**: :beers: separate buy and sell end time, modify cond([`355f0dc`](https://gitlab.tocraw.com/root/toc_trader/commit/355f0dca120df7ffddf14f1e9ffc2c75f16bd23a)) (@TimHsu@M1BP-20210907)
- **debugee**: :twisted_rightwards_arrows: add debug configuration([`d779253`](https://gitlab.tocraw.com/root/toc_trader/commit/d7792539dbd4144f92a36f0b25cf2baaef58aff3)) (@TimHsu@M1BP-20210907)
- **exec**: :twisted_rightwards_arrows: add windows support, change to prompt, split buy sell wait time([`b7f1a0a`](https://gitlab.tocraw.com/root/toc_trader/commit/b7f1a0a277d8f6c7983ff9e168b9736ee4d3c185)) (@TimHsu@M1BP-20210907)
- **kbar**: :alien: add kbar, check target exist in new method, tmp simulate on windows([`551d6fb`](https://gitlab.tocraw.com/root/toc_trader/commit/551d6fbfa3ba4ae8d19337a5f4f9b07a6a774b14)) (@TimHsu@M1BP-20210907)
- **target**: :zap: change target find method, add don't discard over time trade([`7a09a3c`](https://gitlab.tocraw.com/root/toc_trader/commit/7a09a3c69a16c1e4bd47975b07663d2caf20ca07)) (@TimHsu@M1BP-20210907)
- **simulation**: :pencil2: add multi trade day simulate, modifiy fetch entire tick method([`18330ab`](https://gitlab.tocraw.com/root/toc_trader/commit/18330ab1279fda2cc1cb161b228607d6228d656f)) (@TimHsu@M1BP-20210907)
- **api**: :bento: modify trade cond switch api, add manual buy later api([`a195312`](https://gitlab.tocraw.com/root/toc_trader/commit/a195312c5bcd2f7b8309bfcb9f15ec01b8f69a1d)) (@TimHsu@M1BP-20210907)

### Bugs fixed

- **tradebot**: :wastebasket: modify simtrade volume filter, remove check order status == 4([`051f9c7`](https://gitlab.tocraw.com/root/toc_trader/commit/051f9c72fca3b9a1bd25dff95be8e4bb34e14b65)) (@TimHsu@M1BP-20210907)
- **tradebot**: :card_file_box: fix continue cancel same order, no goroutine in middle bot([`2e6a90e`](https://gitlab.tocraw.com/root/toc_trader/commit/2e6a90e1bc9ea1d1bc27945435d1e283a8c5b48e)) (@TimHsu@M1BP-20210907)
- **tradebot**: :globe_with_meridians: fix cancel order flow, modify mutex map lock([`99e9a0e`](https://gitlab.tocraw.com/root/toc_trader/commit/99e9a0e8601361b4ae47d44537516613585220ee)) (@TimHsu@M1BP-20210907)
- **healthcheck**: :tada: remove token logger([`4060a74`](https://gitlab.tocraw.com/root/toc_trader/commit/4060a740221cba4383a8dc33427656d58e88f792)) (@TimHsu@M1BP-20210907)
- **tradebot**: :boom: fix critical wrong delete tradein map, change simulation prompt position([`1d57145`](https://gitlab.tocraw.com/root/toc_trader/commit/1d57145eac8d510a22869d10015b52ae3198f019)) (@TimHsu@M1BP-20210907)
- **tradebot**: :see_no_evil: add closechange in get sell, buylater price, modify cond([`af5f54d`](https://gitlab.tocraw.com/root/toc_trader/commit/af5f54d5f56918452bc53b4dac782577a4d3b7c2)) (@TimHsu@M1BP-20210907)
- **tradebot**: :page_facing_up: fix cancel, add already, none cancel from sinopac_srv([`8727aea`](https://gitlab.tocraw.com/root/toc_trader/commit/8727aead78a477c850c1f025a56ab314c4e0302f)) (@TimHsu@M1BP-20210907)
- **tradebot**: :sparkles: separate tickprocess to forward, reverse, try fix cancel fail but sucess([`8beb028`](https://gitlab.tocraw.com/root/toc_trader/commit/8beb0284e7a708d1c04ea8f497143c67ebfc307c)) (@TimHsu@M1BP-20210907)
- **docker**: :pushpin: fix wrong config path again([`bcf9683`](https://gitlab.tocraw.com/root/toc_trader/commit/bcf9683f574aa3e58210a6a7e5ea0cbe437e900e)) (@TimHsu@M1BP-20210907)
- **docker**: :rocket: fix wrong config path([`f6222bd`](https://gitlab.tocraw.com/root/toc_trader/commit/f6222bd42eba2f7245ac12005c1377075b7b6da1)) (@TimHsu@M1BP-20210907)
- **target**: :white_check_mark: fix wrong target date([`0b06d83`](https://gitlab.tocraw.com/root/toc_trader/commit/0b06d83d0755807d41471ac51c9650fb80b2c0f0)) (@TimHsu@M1BP-20210907)
- **target**: :fire: fix no check day trade([`b86e111`](https://gitlab.tocraw.com/root/toc_trader/commit/b86e1111784b510ebcea052d81fef7a2a60191b6)) (@TimHsu@M1BP-20210907)
- **status**: :package: make sure first status back, calculate ticks per second([`8e5b64c`](https://gitlab.tocraw.com/root/toc_trader/commit/8e5b64c4cc043f48b11685d3a3a672394673830d)) (@TimHsu@M1BP-20210907)
- **stockclose**: :necktie: make sure update close success otherwise fullrestart([`d4c97bb`](https://gitlab.tocraw.com/root/toc_trader/commit/d4c97bbc799b0e99a1ec6854b0f0c9a2c23d84c7)) (@TimHsu@M1BP-20210907)
- **restart**: :chart_with_upwards_trend: change restart method test, change golang to 1.17.2([`cf0252e`](https://gitlab.tocraw.com/root/toc_trader/commit/cf0252e3a6e6fe60d6830292afc68830acf4e47d)) (@TimHsu@M1BP-20210907)

## v1.1.0 (2021-10-12)

### New feature

- **simulate**: :bookmark: disable switch by TSE001, parallel simulation([`a436f88`](https://gitlab.tocraw.com/root/toc_trader/commit/a436f88242a71bdd1b7b9856ba73af346fe45011)) (@TimHsu@M1BP-20210907)
- **simulate**: :white_check_mark: add estimate time, modify cond([`3fd4290`](https://gitlab.tocraw.com/root/toc_trader/commit/3fd429013d4ec9ab3b6eef64473b07c147a9ae56)) (@TimHsu@M1BP-20210907)
- **simulate**: :twisted_rightwards_arrows: save simulation data, improve simulate performance([`0d0dfab`](https://gitlab.tocraw.com/root/toc_trader/commit/0d0dfab1ab92a030a80076f86396dfad7cb8be3e)) (@TimHsu@M1BP-20210907)
- **tradebot**: :see_no_evil: add sell first([`e355335`](https://gitlab.tocraw.com/root/toc_trader/commit/e355335c6b523722099d5a6798c7ba1c1878a89a)) (@TimHsu@M1BP-20210907)

### Bugs fixed

- **tradebot**: :zap: fix short stock bug, re-simulate, add short stock switch by TSE([`2857d29`](https://gitlab.tocraw.com/root/toc_trader/commit/2857d29c0a8c4de7cefb2eaafb9ec45292302319)) (@TimHsu@M1BP-20210907)

## v1.0.0 (2021-10-03)

### New feature

- **simulation**: :arrow_down: add rsi gap to simulation([`b02e75d`](https://gitlab.tocraw.com/root/toc_trader/commit/b02e75da87d57b472d20c0ab65020315205670b6)) (@TimHsu@M1BP-20210907)
- **simulation**: :triangular_flag_on_post: add ticks period to simulation([`957bf81`](https://gitlab.tocraw.com/root/toc_trader/commit/957bf815751f4d1e0aa40427dfe3b6aad14aad98)) (@TimHsu@M1BP-20210907)
- **tradebot**: :zap: change quota method split fee discount, add toggle for simulate, change sell method([`f11a0a8`](https://gitlab.tocraw.com/root/toc_trader/commit/f11a0a8dfd025948910d0df414ec55b7953fbf3e)) (@TimHsu@M1BP-20210907)
- **ci**: :dizzy: artifacts expose as docker logs([`e7b430b`](https://gitlab.tocraw.com/root/toc_trader/commit/e7b430bccd95e74de357b907f9f644ceb36904ea)) (@TimHsu@M1BP-20210907)
- **ci**: :pencil2: add docker log as artifact([`98dc7d1`](https://gitlab.tocraw.com/root/toc_trader/commit/98dc7d1bc569f3d632fef87b648aab99b71894ad)) (@TimHsu@M1BP-20210907)
- **simulate**: :card_file_box: change method to simulate much more cond([`6417c54`](https://gitlab.tocraw.com/root/toc_trader/commit/6417c5412f73acc6a46c0ab57a57a346d314ff57)) (@TimHsu@M1BP-20210907)
- **simulate**: :boom: limit history count to 1300, restart time to 8:15, fix show TSE001([`40d1955`](https://gitlab.tocraw.com/root/toc_trader/commit/40d195570dec91f9ef3e8326cc76005ea438c46b)) (@TimHsu@M1BP-20210907)
- **main**: :construction_worker: add restar full service, modify all api fail message([`9b44704`](https://gitlab.tocraw.com/root/toc_trader/commit/9b44704f6d11ad7114237fcdfaa5260afc9ad3f3)) (@TimHsu@M1BP-20210907)
- **tradebot**: :pencil2: add tse001 snapshot([`f70bddb`](https://gitlab.tocraw.com/root/toc_trader/commit/f70bddb5a0b154d5ed6e60f77dd672b224a263f6)) (@TimHsu@M1BP-20210907)

### Bugs fixed

- **tradebot**: :poop: fix "Trun enable buy off" show even if it is off([`74d8e51`](https://gitlab.tocraw.com/root/toc_trader/commit/74d8e5159ad57438c0f47792fc1b533a1a27cdfa)) (@TimHsu@M1BP-20210907)
- **ci**: :bulb: remove expose as([`bbdaaa6`](https://gitlab.tocraw.com/root/toc_trader/commit/bbdaaa69914519f754fbde2847d8aaa59ee08d60)) (@TimHsu@M1BP-20210907)
- **ci**: :zap: try expose docker logs([`c48b7b0`](https://gitlab.tocraw.com/root/toc_trader/commit/c48b7b0577d889244d93116330f8b8139ca150b5)) (@TimHsu@M1BP-20210907)
- **ci**: :wrench: when docker color log is disable([`b1bf3a8`](https://gitlab.tocraw.com/root/toc_trader/commit/b1bf3a8d8d0810666404d44c869f75193d533463)) (@TimHsu@M1BP-20210907)
- **quota**: :wheelchair: fix wrong stock price([`b31ffbc`](https://gitlab.tocraw.com/root/toc_trader/commit/b31ffbcb01e4a6ecb7a7859ecebb3eab30ccea6f)) (@TimHsu@M1BP-20210907)
- **fullrestart**: :heavy_minus_sign: fix wrong cron, add default key([`5a893c2`](https://gitlab.tocraw.com/root/toc_trader/commit/5a893c21092702852c96a0deaa6d307026fafcb6)) (@TimHsu@M1BP-20210907)
- **tradebot**: :mute: fix TickAnalyzeCondition([`4057f24`](https://gitlab.tocraw.com/root/toc_trader/commit/4057f24e4538086e80ad5c58e0a09526d49dce69)) (@TimHsu@M1BP-20210907)
