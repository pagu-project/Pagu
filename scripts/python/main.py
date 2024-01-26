import csv
import re
import base58
import grpc
import sys
import json

import blockchain_pb2
import blockchain_pb2_grpc

campaign_1 = {
"336952193431240706": ("tpc1pneamej8egf0vlply7cvhszrq0fa48qr93vktux",	  96,	"Adorid | SGTstake#1293"),
"733632046919843880": ("tpc1pfsrcvzf5ce2yuf4qwmzr3g4jasmf5cmzrlad6q",	  105,	"aboka#2166"),
"517631578927661059": ("tpc1p9tq386lawvtt3nks65z8r2yc33ptymyfemhyqy",	  93,	"M0#8364"),
"883636247833219092": ("tpc1pln0wyms36c0qmwhr28z0m0czrh0y2mdg6mtl37",	  101,	"t2N#4911"),
"516728435901726736": ("tpc1pvp5sjr8uwwsnfym97820ps5lpm7uz7w0jg52jp",	  103,	"MeTi#1245"),
"464840045224919041": ("tpc1pa3p0sdm2um4ptsdmjdxhctgq7y66e49kvr5rp2",	  96,	"WebDev | Stake.Works#6225"),
"939985129034637413": ("tpc1psjc0vsu829szdua5vt6w3d57vmk7nht737s5ke",	  104,	"mehrdad kashi#8567"),
"1019299677151186974": ("tpc1pwdykak59yndj5fg3t4c64fxu9yp0vcxhttszg7",	  95,	"BlockDude#1651"),
"998240737592365198": ("tpc1pwn0pcyv57vuh9lhspzvqznuua5gjapjq5t4a3l",	  102,	"Agus Kresna#6018"),
"1007115267287040114": ("tpc1ppqkx93afx54999p7z0zf0lftrm0c90rajaur84",	  101,	"OngTrong#0684"),
"222829790036492289": ("tpc1pv2vn5f45jymdrg5qk8cdvrrvmpjdpjs2shr7fa",	  92,	"BobNymous#7473"),
"860317531470561310": ("tpc1p9q75r2mhqcz36fxnv7mgr6gmf66tujw7czl7gz",	  91,	"Mr HoDL#7879"),
"767249195630329886": ("tpc1pga3ld8uq9523sss87tdwtrvhy50dgvn5tvzsfq",	  102,	"jwmdev#3991"),
"948574466743607358": ("tpc1pwl38jqsqwlcx9u64yng6kv4gvwzgnwxdp2j2uh",	  101,	"Yanz#2600"),
"934015758785212436": ("tpc1p025lrf2e56u6ky229yt7lt75xwy2u2k8chqnru",	  91,	"genznodes"),
"1081629757793374218": ("tpc1p37wknz8zfenpwl2psksp9ss8dtrnm6p02pkayk",	  91,	"kehiy"),
"994159473558040587": ("tpc1p7lftg8q5n96l3kajcdke3ny5n9p3aka8p66z2j",	  91,	"nn0ki"),
"1053059894347038870": ("tpc1p0a87plfqu77nzj8rccumjr3wg8jc0ckmsw789l",	  101,	"crezeta"),
"1052589576092401744": ("tpc1p2vs07lh0md2kwvg094xd2k2shpkzncjrqzqd84",	  101,	"shdxid"),
"957351281427636284": ("tpc1pk09xnkrfy47uwvn4prqfq6f29lxsa8lxjetn02",	  91,	"g4rver"),
"1050786675208507462": ("tpc1p8dddx9vfhxelltkx6pf3d4pjewv084dtjwu4rr",	  91,	"gvnhz"),
"827586639366717472": ("tpc1pj3k4pv5ujttd9gh943q8yf8xduaq0l79u00pnc",	  91,	"kv1ar"),
"957413301158047814": ("tpc1p8wqgmagsrzn0nr26weg6wekqtu2mc6uw72k04a",	  91,	"at3rnd"),
"705747180887212104": ("tpc1pua7n97uftsehmuynlhye7nkdz7f2q9hq9rttf8",	  91,	"dfxsss"),
"356290864340926464": ("tpc1pxpmrnjn6wn9upjwu5ee223z6rvm76fnfx3u4l7",	  91,	"etherscan.io"),
"910185403385020488": ("tpc1puc5zza3hnp2tcf6r5n8zz0mwcjhqlxtejnjkzv",	  91,	"warriorcarl"),
"997388215755477163": ("tpc1pj0k7zpthh82tl394cz3pns5tthrc8fzqvfaupp",	  91,	"april#0537"),
"842691688656666656": ("tpc1ptdsqzwvq5h4tmqy3vgmu7eqadyr65epxrlr396",	  91,	"jackytbe#3999"),
"213018208079183872": ("tpc1p5mwy7tfdva2e9z736gjsftwjgtrnv3ujf8jw6w",	  91,	"@0956"),
"479237981610442762": ("tpc1p5ze7r3q3m6k60lhpqxfsg854840h8q3yyyv5qs",	  91,	"wzsd"),
"837304444014034944": ("tpc1pt5gy4s32q5aywq823nry9dgffszz2zx80z4n4n",	  91,	"faturalhusni"),
"783201902807875636": ("tpc1p27t6pj4r736034vapspw2pkualvtj823pr26af",	  91,	"9oal"),
"948854452272652318": ("tpc1pwdxdavlauqu2073ulpjhv2zf3pq9m98d7exj7r",	  91,	"sledgerhammer"),
"773184173355565056": ("tpc1pj27rf4e96h4rn0xhdcxp8nmvuz27t0ffe0calk",	  91,	"cypher_knight_007"),
"932167886221504573": ("tpc1psqh3g8q267py87re9d92gr67eykzvtwm8zq493",	  91,	"xasla"),
"981397790947180545": ("tpc1ptmtrze38exrmp5d33ck6twzkcny3ayqfejwvj4",	  91,	"fachrulspc"),
"399571183986802688": ("tpc1pupr676rvylppvzyy36rr9chjhgkphmjfwtzwgf",	  91,	"catsmile91#4043"),
"1041933478733828146": ("tpc1pgup9ud0kfmkun0qxhdch7xp6f4x963ujltjdw9",	  91,	"nara.web3"),
"927161736837083157": ("tpc1p73wj5l48lpm5wetpazpczg3peqkxm28p42jrmw",	  91,	"abhi#3886"),
"799242450186665984": ("tpc1p5mtrvt6ga2uvy3xyyyeghxxjfxcavmxqm8dret",	  91,	"ovzx"),
"908179770611728425": ("tpc1pectlfgufn52atyvg9shffk6qn588xftzwmswhf",	  91,	"gallerynft"),
"426773473789214730": ("tpc1pfpvuxumq59uknw5jfmpqrexhl5ryva2jl4t88p",	  91,	"peellygg"),
"678931901427482656": ("tpc1psne8ypas22nfue2hwgl8n48ppzz3f2rx6lahw2",	  91,	"0xRyuuki#0"),
"384758639946235916": ("tpc1pmaqteanyd8cgwfpzjwljpf3f744dj86n0veuh7",	  91,	"sheza_74"),
"907280152910782514": ("tpc1pgw8tyxxkxykwxt63ecwhgmnq2am9nnaumv6smq",	  91,	".inferno46"),
"831214245664129024": ("tpc1pc48cxvs87g2n5nwjnwru4y92q8qpzellq0c2l6",	  91,	"zhuxan#6636"),
"953682756116840468": ("tpc1pv4fxln6xec6tu7dmvml6ykjxvqlf5ff67qqmyc",	  91,	"arbitrumdao"),
"905802225287303218": ("tpc1pjc6v92j0mqhxj2qlrf44u0y7j9pnhm733aaska",	  91,	"@hnf1407"),
"896380682513829909": ("tpc1pn50sqxas7xa3acunk337wwvu9x4mks2xzlfk2f",	  91,	"diabloo#4791"),
"1113904869493985280": ("tpc1pplc30rj929m4r586utu7meurvyugjqu7cr9qe7",	  91,	"kiddstark9"),
"641663460119150602": ("tpc1p6t9axlc92mmdnz2vzpw3nkrzmmh9yjjrsgsw4u",	  91,	"zcode"),
"455742059379425291": ("tpc1pez9ptc3c84zekk0d230jxaaykcpgx0k70tecj7",	  91,	"goneth"),
"961673869410852924": ("tpc1p6a0e4ywrnyggg6sfhz546nn0rk079j9ka9wra0",	  91,	"arieferdieansyah"),
"890905491201490944": ("tpc1pvx3zxk2pz3kucpyapzlvp5swr5g0vys48cs9qk",	  91,	"jordialter"),
"1109620728904568974": ("tpc1pldh8ewlh6rmgghnynpx5zqlt5agk3n7s3yaymg",	  91,	"ai.93"),
"1023600450223755395": ("tpc1p82jjr05rwwj7yvx6jc7fk5qsv9gf5atq0qy2jj",	  91,	"sideko94"),
"425313215610617858": ("tpc1p0dlrwq94fcz9klxz3rk45mx3jqy33qqtpcf306",	  91,	"pa3l"),
"579525123644588033": ("tpc1p5ckqv09slkcmj2tet9w3ta997yktzdzdc3gep4",	  91,	"21.btc"),
"882910372460388362": ("tpc1phv4rwpx7nwxpfzt2dwr42d8q4demwgxhnm8dgx",	  91,	"frankie0042"),
"885630036563091486": ("tpc1p9njmv20dpxd8a5hl9y7tt5fder6gvtgp6xet6l",	  91,	"Nae#1271"),
"925389883416141877": ("tpc1pqjc4am9pvhxc3wqunc7k9epfrnukqld2586c3x",	  91,	"mpola"),
"907044201345191946": ("tpc1p57n7ddqvxpc3dchpqtnck9nrrzr7vd3glmurgk",	  91,	"Lilik3004"),
"1060746133460242492": ("tpc1pq7tf6hlpqsrh4lpczcpa0qh4pf72nqlrrd62d6",	  91,	"mymass"),
"1091787351212179526": ("tpc1ps82jwh4avljx7rgqtjf0vea7ekckahhkyckrtp",	  91,	"suryor#8182"),
"1065780683471081502": ("tpc1pvu5wf7u4r0elet0dns4lc4jf0rv7x98hp3xlhy",	  91,	"bensol1904#2166"),
"1026063705286377482": ("tpc1p7wfem69lgewhyl990as5e9xn22t777g3swvjs2",	  91,	"AISS#4411"),
"462222091522146305": ("tpc1phepzp076x2teu52nxtujkcgylcd3gvsqv7vvvm",	  91,	"g3mbok"),
"820395564272713739": ("tpc1pa9qdjpqaqhggd2wulzjqlfqvkgqwy8dlv70780",	  91,	"sunewbie"),
"950492754075598898": ("tpc1plqvy3at8pmat7e3jcggjeqc28vhhh4dq4eqqtd",	  91,	"chikabul#0"),
"794806355088637953": ("tpc1pajjh8d9fw7jlgvggx9u050arvkxssynfuwt2lt",	  91,	"! OX3cDF"),
"1002269095170945094": ("tpc1pxeucym9jdlzusqjvd0x7hntj5k78n9rc2l2duy",	  91,	"Jaboeybae#7550"),
"547986237147971584": ("tpc1pcpq60fg3eyr990w5pgzg9a0ljksxvyd7cxxqyj",	  91,	"hendrazlk#6590"),
"841961575685554208": ("tpc1pu89tdj9x72gwgeaellqtsyxmgxs2ejv62l8fqm",	  91,	"batex_o"),
"766587388528558081": ("tpc1p0rdtmuxqcw22taa7ts3k2en896x9r46cmwrcff",	  91,	"amongussus#3448"),
"488896313279119370": ("tpc1pfdxhxnwf46qmy0tmnadrvgfchayuh9ld4wv4uc",	  91,	"gofur_triad#8493"),
"856542400464551987": ("tpc1p83vdgz7u0lc822rdy8yq2zg8c80uadf8m6ll8q",	  91,	"zulkarham"),
"944157601388724234": ("tpc1p60vch4hmkyztvxkeal4caqqegwamagvxu2q8ah",	  91,	"indahnuralifah"),
"984060206155710484": ("tpc1pd49qarx4atq066ld5qt6tu5s8hdva8nkk8gt96",	  91,	"sipaling"),
"444109910003941376": ("tpc1pgh5d6lz9zq66mquzvsslfv28x4wghny9w50zu2",	  91,	"malghz"),
"1041061690684489730": ("tpc1pepqck94ln3p8v00p934l5p7lsy9ul0k7c5ts4k",	  91,	"mashiaplghhg_759"),
"949725569682133072": ("tpc1plrdv9gpsc6x266qr92a2kp99d8qenvgj7qm3ff",	  91,	"megga#8040"),
"577300640095535124": ("tpc1pcz7lt08rg7m38sj7ne9srxnme8x9p43ysqsvpz",	  91,	"0xRgp#1618"),
"948088507849642024": ("tpc1pvlc4lv8uteva3l9mpe9rdedmv0v2fg43hetmzj",	  91,	"akumantanmu#2952"),
"764142104972623913": ("tpc1ph0kq87wedpd8u6ms5d7pke8dljfvcuvhaykp9u",	  91,	"keqingwangy"),
"1071998655634100376": ("tpc1p4r5jxcjtskdzzq2zxg900gtnfpjgdtvnx68gak",	  91,	"pheromone#1040"),
"835731103171608627": ("tpc1pwlett3qjy5p4xwdyqmw4v0ysesggg3mtwl3xmg",	  91,	"anggawrt"),
"1104932144431775755": ("tpc1pta0vtwvvd2cw7rrf49eyw6qwm2nrjh2eghr7cx",	  91,	"solehanam1#3458"),
"1108378542498140240": ("tpc1pf8qappg4evgp6xznznw2enlx89zf0vqdw3l74h",	  91,	"cyfan100#9777"),
"848078128123346945": ("tpc1p78m5u4nc2t8ks0j5ltvks8pduteqfprdl7nycp",	  91,	"nengRahma#2670"),
"841826423441457152": ("tpc1pz4zszxqnxr0ymmtu0m68fu00wuwccjqkxa2v8m",	  91,	"homeduoc#5515"),
"1094139822089703454": ("tpc1p8pe7f6dn2qqc37xuqta5t4fdu9nsddu9ymxn8m",	  91,	"lisamiran"),
"424422625662468096": ("tpc1p3vdkq58cm04pxh7fnc8l6453kfdnrx463ux7at",	  91,	".syanodes"),
"1117314911484256386": ("tpc1p6jtg4cct6s3kzdh6t3qzcnq7ue9ju72eat7yyk",	  91,	"thaitokenlegends#3619"),
"906483432811561000": ("tpc1pcj8rp29nudfgp0sh9ng33xu9vtrmx4xhnrjndl",	  91,	".shazoe"),
"960679963386871848": ("tpc1p32yzazm58mfzhrenyz67v9gn3dnzfxpzkan2rf",	  91,	"remix.ethereum"),
"445212864815562754": ("tpc1pvuhnfjne9tel5nsu7dvytha6spugcu53pslug4",	  91,	"saandy"),
"342239807876890625": ("tpc1pwdzlrcuk70l5hp0w6ykfc53cxe9carnz5wcyts",	  91,	"djieyz#6051"),
"932572788903002152": ("tpc1pqqvj9nuxlqm0plheekze49fwtx5n5pvttpw2zk",	  91,	"qarambytre#9532"),
"773407544340381697": ("tpc1p3xxmk0wrmzv09cf8ej8j9j9eclvs8rzul5yzg8",	  91,	"aidil.sol"),
"917275641441816607": ("tpc1pmf879a95s293h77885kg4vk5zh46k35q4fu6gm",	  91,	"babangaip#7848"),
"221664017859608576": ("tpc1p2zwmydmrtl5rm7g2e0jmk328gj9j0djkwukfp2",	  91,	"zianlin"),
"422805009265197067": ("tpc1p2hqzc999uenkrnd4jhngkrwxlrc763fssezayg",	  91,	"Logosdibta#1882"),
"841628058934575114": ("tpc1psze0jul5ggpq36cnyhm7l3tgrvv9fquxj62prh",	  91,	"strnan"),
"896408959903207484": ("tpc1p02qpv3fmc7dts8n3h2j8wn22tct8crc8txrc9q",	  91,	"Caffein#4863"),
"357698971445231627": ("tpc1pl4eyv8n3krs7t74apxax8n7gpxgfde0au240pw",	  91,	"pramonoutomo"),
"1056208697468125194": ("tpc1p6pr8s8rzqs9t6enf57v6zenrh6hduauje78j6y",	  91,	"James77#7719"),
"571031781718097920": ("tpc1p8sa9c6f5krwgp985gz8ezm0e9yx60gqexh7fps",	  91,	"@vermillionss"),
"858626777327730718": ("tpc1pc00hkm0xc6et083w2vyma2px833fz3vfu430cm",	  91,	"thuongtin162002#6399"),
"853852856557633561": ("tpc1p9qhvzjr3q3qlgy5y50nhj7lm28uqzv97qnuhnk",	  91,	"tinboy_"),
"907302551383318548": ("tpc1pqk556c0d4ygr76spgl3xq7hfjmt5kga044etlv",	  91,	".boester"),
"1035118911324172298": ("tpc1pk8nlfzupy5gjs2p22kcy46gcxcy7l0d9up4fvy",	  91,	"khangnguyen3790"),
"457941630834442241": ("tpc1pkyavh2amx53sxaaj4qpfhzpkg0gxxytyeyl2yq",	  10,	"mahendra_"),
"494313350209994765": ("tpc1p06jxcgqrvhdtp58fl9ua69a8n0drm4qfp0sdl9",	  10,	"irianty46|Beagleswap"),
"865254852644175893": ("tpc1pzj4lmseadwrjg5wjv8xech54l3msaep55k3ma6",	  10,	"up.bit"),
"1061684503082446948": ("tpc1pv5jnczsjc95rzc2hy83udx5lrfrmev5zwc6w5s",	  10,	"radyspradana"),
"788336120056250428": ("tpc1p3w5rc6yuyxjpe0tufvlt3f4237c704w79dat6u",	  10,	"@shunna05"),
"1004277158069403659": ("tpc1p5lgah9hqeeqtjsuce8ljdzw3tczca68w6eq9p7",	  10,	"Aileen#1743"),
"515154342395772948": ("tpc1p823x9enwqwsxpjt6g36lqskz74c3vuc2e8a2j4",	  10,	"morz#8861"),
"848060688852189207": ("tpc1p0t7qw89yeh6e4psph67zejdcfktqdyfuuhdrht",	  10,	"pathum2223"),
"1073090864454307880": ("tpc1paa3wjvht5v2y86lg08f27rre7uhdwcf6vgyxp0",	  10,	"asyajoker#8856"),
"1011764278702907412": ("tpc1pagq8wlmhe8r0d0uknvuv4q3q434ps9cznaeag9",	  10,	"artra97"),
"883693707105292359": ("tpc1p97n8vfguwldpef682lauzj6czvwehts5y6q7xe",	  10,	"atalasia#6472"),
"1069537067706630174": ("tpc1pyqxeq6n9n2qcc86trg0t7upzvyy92xvtkfn6lk",	  10,	"blacan#6803"),
"509529719671095327": ("tpc1pdy6xy4d3f2nrq9wl0zlsv6ht5ka9qjf3ykacyw",	  10,	"0xrizal"),
"324780560138371072": ("tpc1pzsrvc88cexrnag7g2rgjt2f4qquge9dyemtcek",	  10,	"Invalid_validator"),
"893931014999658496": ("tpc1php0s9t7e0w2nvn7g6p3r4a5wj2auu6e4ckqh5x",	  10,	"laedrei"),
"927171117511225435": ("tpc1prz3mstqdes0nhmm9z29sae7d6juar23wkdhx6t",	  10,	"0xl_"),
}

banned_list = [
    # "tpc1pfkk7zhc3xmnxkf22nqx4s76dsqzqsgzyctggng",
    "tpc1p60q9y3639e4rjnwqe2rcwl0m5nygca2zpv9p3g",
    "tpc1pphezl7d84k25wftsdn8498tcjhw3xrumpkmzzq",
    ######
    "tpc1pas2tv2tkm3gqdxytrqespgtg8udfsqk5ky9k0k",
    # "tpc1psv7cezxnu3lt59z4qd3pparwyhytu5fnpd6qfl",
    "tpc1pyk0ecqtxt9aakl0n7spvh35wgnez9ruf0twf68",
    "tpc1pdzljcwt0zje9zks6kmvslyukrjpc6v0642hnk7",
    "tpc1patefjqx6gsqzvpan7vyg4ceqz3fmdfnwwpt2g2",
    ######
    # "tpc1p3uq7l8uny2uxut6wkkd2nnf4am2egmfx4l08ae",
    "tpc1pev3ncqnwqfu0ak05h2sdf6ptz40n9m6stp55g8",
    ######
    "tpc1pkdf49w7eqalhhk4yp2d9p9sd4wr7wz2h83pukz",
    "tpc1p8sxun2wx4zwa9vs30c0tanlyu2rg56ylwc27hd",
    "tpc1p3pf3ue0g62gsaykszf09nntl4fcdqfq0zqf2n0",
    "tpc1pmupwwmn7gxjutlgrw2j8jler68jzfv0ywx6djv",
    "tpc1pujmqcc9h0n5n7zru49ftvlgwz5rpggxeenh5m2",
    "tpc1pf26qydcgz5nt7rxn6qljkfyvjvywqcgxhsy0f2",
    "tpc1ptdumpgtz3xdmg2q9d8esk8lwp46ctlrns69m5l",
    # "tpc1pp28xy9mefmtl6rxghvvkr3kwwknacvpau42ey3",
    "tpc1p5sz43z7ps3ezj0dpnca5g47f587xtcdjwvcyd5",
    "tpc1ptg5euxr6m5dcllmxghcc8suuch4x8rn03ky42n",
    ######
    # "tpc1pfmmwr7k2rppj09khgqsj3y48m0euald948ke3e",
    "tpc1pvutf2fmrhw6jkffzrf5r7tmztusg436f02qffc",
    ######
    "tpc1p0m4yqchwyx96sygelfuhdwm5tcwx39ffwqmen4",
    # "tpc1pt6x5053p6n8gau0nu6ahsrcyz0p5ngxc72n682",
    ######
    "tpc1ptzfu5s7uz0350rr4tek0sxjjd7tlcyvd7r2w67",
    "tpc1peajdmvjqqxyh866e389qas7dzu06aaqx7r7kaf",
    "tpc1pmlx38mljvnnu95a470tzashwx4rujwef6afsvn",
    "tpc1p66xfdrxmm4pt9grqqydh3jr335hvgcxhjm52mn",
    "tpc1pa4mukkc46uk55en7ruezmtf0ce7ur4dr4kd3dq",
    "tpc1pgfxw2hn9my988zgmsgw45d3zlxkerluqk8p06j",
    # "tpc1pudh2el8cszep5p0x6vfentfqm49hhe8c88jjxe",
    "tpc1pdtjqrv0ktxxhnneyhfjypfl7kpy4ecu4jznje2",
    "tpc1pdfrr4j4kxqxswefct59wkfyzp5amycgylv56ln",
    "tpc1p6u9xvwlltfgjdaqv6vxj6ragmed76mvuw4mgqm",
    "tpc1pv4l4dkh5v48nzftpxcl4lelnwdq2jm36ymc6yj",
    "tpc1phgn2hrvk3j3w35a8df7zl39xpmq7ml0gduu0dj",
    "tpc1pndhgc070evjnmymdnk5vm3z83a5wslextekfpe",
    "tpc1pnd2ut3jrasx3nwj03j4y5nrrv3aadwehvr09nc",
    ######
    "tpc1p7lt8ylswzzzgge4wlsmrkp3esspgs6v8dnzev7",
    # "tpc1payqvcpulrg7rdnq3jg3wkc5w0z88039syzl7ry",
    ######
    "tpc1pwvqd8s3xtd60jemp2lsevqwxg5wqnv0u8zcn23",
    # "tpc1pmp264udls6sq5nr0pre3ted8j2f7nf72xy2d5l",
    ######
    "tpc1pl7l4aw80jpyt9dtngleeu0a7lw4s5hceetpksz",
    # "tpc1pvdd4nrtlcfsykw2yrt666zm4xyztle05svevsf",
    "tpc1pcqpw4v62c6skedr49hk3q8aj4hmhthd7l9d8ff",
    "tpc1p5hjk6ht7wyv66raqf4hj0jlg6s4k7t06uw02r9",
    ######
    # "tpc1pva4zxc0vy8mzhmdvttms7g04eelrdwapc59w9a",
    "tpc1pgdpkqfvtsjxycyu8l8z7t7jdwesn5eqm8nc3xc",
    ######
    # "tpc1plgczlh8lxpxg2lz67gzkp83rpy2r3yxaua5s5a",
    "tpc1pk93u0efxx3epl7gnfuqakk7p3rjdqpvlhqeufj",
    ######
    "tpc1pqds4hc60nyfd45e3zphcf3n90raggf2yn04dlx",
    # "tpc1pnnh9a84zuvgnzflluexkt57nxd4jel7mvp6yc7",
    "tpc1pn8dnvny63ap73y2dwxtq9jvshdt2uut6luhv6r",
    ######
    # "tpc1ph25hveeg02zmzmz8pc5e85rs3d993j9lw5jjuv",
    "tpc1pht5hyu9twk0mg9xym9l7au7t9g2h5zj6g8y3dd",
    "tpc1p3fu0l7r5g4jepnulyc4gpj8wd5anrmf5j0ced2",
    ######
    "tpc1pjs6xyexg2gsh4jymzp5pzn5ypzlt8sk760vumz",
    # "tpc1pck7q7eyljhfe7h5ms3vfdtd84gelhuyf0remaj",
    "tpc1pqv4kel45y7s5ltpxw3ya0gcamwn3mw05je9we0",
    ######
    "tpc1pz94nrlld8ay39tejxk2xrxpslwqgrpdelr58hx",
    # "tpc1p20xcx28707m4reh396lymvsn06ffwd2z0rku6d",
    ######
    # "tpc1prqstqp9qsmv7lsp2w8qzvxksxdeu7pnasercc3",
    "tpc1pmjwv2a4uvvg0m3gmj32ya7tay3fz9x898fkexj",
    "tpc1php06vkd9j4nup3t3evhkuqkaw4jqhj5yzefg96",
    "tpc1pl9a2z02mr2mzhhnwhttvpf3m26qsdptzsxekau",
    "tpc1pdeenwahh4puz74txk0vufanuv6j6wulln27zum",
    "tpc1pv7pprerrhfytejkm8shha9k6e9vc0yevdey2xs",
    "tpc1pstuygk6d0xl78pg4gssh2ne02nd8pkkx5j8c4k",
    "tpc1p7gpdctr3aekvur052kdaz9sspcacdn0pw0qfx5",
    "tpc1pazz5eu9tc9hqd070y0vtfdq0xjps059ypqk0j9",
    "tpc1psutqmydgg48frd2egt45wy8zva7fjs60we05gm",
    ######
    "tpc1prfggpj23dmxy70qt8n4zhw6wnr7h4w4xms7kem",
    # "tpc1pspzwgzlu6edyev7lkky27gyjr3dw6zjmuj5hft",
    "tpc1pkg0lc59wyu9us2vk6mxkpj7n7ve9des75td0rw",
    ######
    "tpc1pecl3al2zkd3xe0w66lgwm8ug6ngm4mud9mcxfl",
    # "tpc1pnf5quxvme3xnmfrj6d5zcnfe0w73xkkgypgzkp",
    "tpc1ppa9c0tdm6g6rx22ahhztgkn4kvuyjjpjsdsyqk",
    ######
    "tpc1p2fudhnln424wzrmrsy46ck3sr2h37pndaclaal",
    # "tpc1pcpxu9sqcx82wal4rzqexa3aeadxglkwxfld535",
    ######
    "tpc1pl5arrhjvplt2f8q4t559jz584kpsklzrkkn22p",
    "tpc1pdq4crht67mdvefjlmdhvtwrjm0cqedemt5vdt9",
    "tpc1p3dxjf9lr5m9sgt525kyxzxdutgr02akyjh6v8k",
    # "tpc1pffrcfds5ytn5vj7hfernjg9pwf88uzc37et2s7",
    "tpc1pk0d6gj33my9dvc0fudcnukz79hncr98d3tgvmc",
    "tpc1prpq9xndu4nrjlh6n8fx2sjfg98f4zvguraz3h3",
    "tpc1phx2qx34tq56q4cajepxzttzwc8ekh2vncvfp7u",
    "tpc1plp4l2g8qzsrwx02wafdn7x0qp8gyfhdedckfvh",
    "tpc1pwmwm7nqzxt9ezlurdz8n3uswj5mjls2q3dpam6",
    "tpc1pj7ct4lahk3yq45r6mhqtxc406htjd0hhwugwd4",
    "tpc1phcu59wds3xuczw7axzkzh8xcgap337s75svh0y",
    "tpc1ppyd6rgfl049ze7up2fwks3hcydtujvkwwlc05m",
    ######
    "tpc1pqpqrn3jd66z2w3d822s386dpfua4qgpgvp942v",
    # "tpc1p4sm9grn8rr0hng7cjls3szuxex28ayn0v7y5k9",
    ######
    # "tpc1p3cx5hqh0604whu84ydmajl4nkyf8q739ze2fz6",
    "tpc1pwzvuvaezhks7tl5v3qdg0ntruh7lmyv2v7sw26",
    ######
    # "tpc1pk39tyyhumyxysd50kztj9e04euuyjj9zlrqg3y",
    "tpc1pplxqpupgtts6j09fflswh4tz8t8cdm0u5kke2n",
    "tpc1pa7yruh0t2c8xqv40ht6mr7f9u4qdymmlr56ekm",
]


def load_csv_as_dict(file_path):
    data = {}
    with open(file_path, "r") as csv_file:
        csv_reader = csv.DictReader(csv_file)
        for row in csv_reader:
            data[row["Address"]] = row
    return data


def load_json_file(file_path):
    try:
        with open(file_path, "r") as file:
            data = json.load(file)
        return data
    except FileNotFoundError:
        print(f"File not found: {file_path}")
        return None
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}")
        return None


def get_validator_info(grpc_validator_stub, addr):
    req = blockchain_pb2.GetValidatorRequest(address=addr)
    try:
        res = grpc_validator_stub.GetValidator(req)

        return res
    except:
        return None


def referrals_by_discord_id(referral_data):
    # example of referral data
    #
    #  "039370": {
    #   "referral_code": "039370",
    #   "points": 0,
    #   "discord_name": "xulee99",
    #   "discord_id": "923711288759177277"
    # },
    referrals = {}
    for key, value in referral_data.items():
        referral_code = value["referral_code"]
        points = value["points"]
        discord_name = value["discord_name"]
        discord_id = value["discord_id"]

        referral_data = referrals.get(discord_id, None)
        if referral_data is None:
            if referral_code != key:
                print("Duplicated referral code")
                sys.exit(1)

            referrals[discord_id] = {
                "referral_code": referral_code,
                "points": points,
                "discord_name": discord_name,
                "discord_id": discord_id,
            }
        else:
            print("Duplicated referral")
            sys.exit(1)

    return referrals


def extract_users_map(validator_data, referrals):
    # Example of validator data
    #
    # "12D3KooWSbDJWMhYgrFqymf78q4vhZEAV8n1LUUZ7Y9VenR6PNdN": {
    #   "discord_name": "vikanren",
    #   "discord_id": "840519270004686878",
    #   "validator_address": "tpc1pjrvumvpsutpklgg2hwhuaaujc3pju6z9kkffzt",
    #   "referrer_discord_id": "",
    #   "faucet_amount": 100
    # },
    users_map = {}
    for _, value in validator_data.items():
        discord_name = value["discord_name"]
        discord_id = value["discord_id"]
        validator_address = value["validator_address"]
        referrer_discord_id = value["referrer_discord_id"]
        faucet_amount = value["faucet_amount"]

        user_data = users_map.get(discord_id, None)
        if user_data is None:
            user_data = users_map[discord_id] = {
                "discord_id": discord_id,
                "discord_name": set(),
                "campaign_1": 0,
                "faucet_amount": [],
                "referrer_discord_id": [],
                "referrer_discord_name": [],
                "referral_points": 0,
                "referral_code": 0,
                "stakes": "",
                "total_stakes": 0,
                "total_reward": 0,
                "total_stakes_online": 0,
                "total_reward_online": 0,
                "num_validators": 0,
                "validators": [],
            }

        user_data["discord_name"].add(discord_name)
        user_data["faucet_amount"].append(faucet_amount)
        user_data["validators"].append(validator_address)
        user_data["num_validators"] += 1

        if discord_id in referrals:
            ref_data = referrals[discord_id]
            referral_code = ref_data["referral_code"]

            user_data["referral_points"] = ref_data["points"]
            user_data["referral_code"] = referral_code

        if referrer_discord_id != "":
            ref_data = referrals[referrer_discord_id]

            user_data["referrer_discord_id"].append(referrer_discord_id)
            user_data["referrer_discord_name"].append(ref_data["discord_name"])

    return users_map


def write_users_map(users_map):
    # Specify the CSV file path
    csv_file_path = "output/users_map.csv"

    # Create a list to store the rows of data for the CSV
    csv_data_users = []

    for _, value in users_map.items():
        csv_row = [
            value["discord_id"],
            value["discord_name"],
            value["campaign_1"],
            value["total_stakes_online"],
            value["total_reward_online"],
            value["total_stakes"],
            value["total_reward"],
            value["faucet_amount"],
            value["referrer_discord_id"],
            value["referrer_discord_name"],
            value["referral_points"],
            value["referral_code"],
            value["num_validators"],
            value["validators"],
            value["stakes"],
        ]

        csv_data_users.append(csv_row)

    try:
        with open(csv_file_path, mode="w", newline="") as csv_file:
            csv_writer = csv.writer(csv_file)
            # Write the header row
            header = [
                "Discord ID",
                "Discord Name",
                "Campaign 1",
                "Total Stakes (Online)",
                "Total Rewards (Online)",
                "Total Stakes",
                "Total Rewards",
                "Faucet amount",
                "Referrer Discord Id",
                "Referrer Discord Name",
                "Referral Points",
                "Referral Code",
                "Num of Validators",
                "Validators",
                "Stakes",
            ]
            csv_writer.writerow(header)
            # Write the data rows
            csv_writer.writerows(csv_data_users)
        print(f"Data saved to {csv_file_path}")
    except Exception as e:
        print(f"Error saving data to CSV: {e}")


def extract_vals_map(validator_data, referrals, network_vals):
    # Example of validator data
    #
    # "12D3KooWSbDJWMhYgrFqymf78q4vhZEAV8n1LUUZ7Y9VenR6PNdN": {
    #   "discord_name": "vikanren",
    #   "discord_id": "840519270004686878",
    #   "validator_address": "tpc1pjrvumvpsutpklgg2hwhuaaujc3pju6z9kkffzt",
    #   "referrer_discord_id": "",
    #   "faucet_amount": 100
    # },
    vals_map = {}
    for _, value in validator_data.items():
        discord_name = value["discord_name"]
        discord_id = value["discord_id"]
        validator_address = value["validator_address"]
        referrer_discord_id = value["referrer_discord_id"]
        faucet_amount = value["faucet_amount"]

        val_data = vals_map.get(validator_address, None)
        if val_data is None:
            stake = "0"
            last_received = 0
            last_sortition_height = "0"
            last_bonding_height = "0"
            unbonding_height = "0"
            availability_score = "0"

            if validator_address in network_vals:
                network_val_info = network_vals[validator_address]

                last_received = int(network_val_info["Last Time Online"])
                stake = float(network_val_info["Stake"])
                last_sortition_height = int(network_val_info["Last Sortition Height"])
                last_bonding_height = int(network_val_info["Last Bonding Height"])
                unbonding_height = int(network_val_info["Unbonding Height"])
                availability_score = float(network_val_info["Availability Score"])
            else:
                val_node_info = get_validator_info(
                    grpc_validator_stub, validator_address
                )

                stake = val_node_info.validator.stake / 10**9
                last_sortition_height = val_node_info.validator.last_sortition_height
                last_bonding_height = val_node_info.validator.last_bonding_height
                unbonding_height = val_node_info.validator.unbonding_height
                availability_score = val_node_info.validator.availability_score

            val_data = vals_map[validator_address] = {
                "validator_address": validator_address,
                "discord_id": discord_id,
                "discord_name": discord_name,
                "total_rewards": 0,
                "campaign_1": 0,
                "referral_code": 0,
                "referral_points": 0,
                "faucet_amount": faucet_amount,
                "stake": stake,
                "referrer_discord_id": referrer_discord_id,
                "referrer_discord_name": "",
                "last_received": last_received,
                "last_sortition_height": last_sortition_height,
                "last_bonding_height": last_bonding_height,
                "unbonding_height": unbonding_height,
                "availability_score": availability_score,
            }

        if referrer_discord_id != "":
            ref_data = referrals[referrer_discord_id]

            val_data["referrer_discord_name"] = ref_data["discord_name"]

    return vals_map


def write_vals_map(vals_map):
    # Specify the CSV file path
    csv_file_path = "output/vals_map.csv"

    # Create a list to store the rows of data for the CSV
    csv_data_vals = []

    for _, value in vals_map.items():
        csv_row = [
            value["validator_address"],
            value["discord_id"],
            value["discord_name"],
            value["total_rewards"],
            value["campaign_1"],
            value["referral_code"],
            value["referral_points"],
            value["faucet_amount"],
            value["stake"],
            value["referrer_discord_id"],
            value["referrer_discord_name"],
            value["last_received"],
            value["last_sortition_height"],
            value["last_bonding_height"],
            value["unbonding_height"],
            value["availability_score"],
        ]

        csv_data_vals.append(csv_row)

    try:
        with open(csv_file_path, mode="w", newline="") as csv_file:
            csv_writer = csv.writer(csv_file)
            # Write the header row
            header = [
                "Validator Address",
                "Discord Id",
                "Discord Name",
                "Total Rewards",
                "Campaign 1",
                "Referral Code",
                "Referral Pints",
                "Faucet Amount",
                "Stake",
                "Referrer Discord Id",
                "Referrer Discord Name",
                "Last Time Online",
                "Last Sortition Height",
                "Last Bonding Height",
                "Unbonding Height",
                "Availability Score",
            ]
            csv_writer.writerow(header)
            # Write the data rows
            csv_writer.writerows(csv_data_vals)
        print(f"Data saved to {csv_file_path}")
    except Exception as e:
        print(f"Error saving data to CSV: {e}")


def apply_campaign_1_to_users(users_map):
    for discord_id, data in campaign_1.items():
        user_data = users_map.get(discord_id, None)
        if user_data is None:
            user_data = users_map[discord_id] = {
                "discord_id": discord_id,
                "discord_name": {data[2]},
                "campaign_1": 0,
                "faucet_amount": [],
                "referrer_discord_id": [],
                "referrer_discord_name": [],
                "referral_points": 0,
                "referral_code": 0,
                "stakes": "",
                "total_stakes": 0,
                "total_reward": 0,
                "total_stakes_online": 0,
                "total_reward_online": 0,
                "num_validators": 1,
                "validators": [data[0]],
            }

        user_data["campaign_1"] = data[1]


def apply_campaign_1_to_vals(vals_map):
    for campaign_1_discord_id, data in campaign_1.items():
        found = False
        for _, value in vals_map.items():
            campaign_2_discord_id = value["discord_id"]

            if campaign_1_discord_id == campaign_2_discord_id:
                value["campaign_1"] = data[1]
                found = True
                break

        if found == False:
            vals_map[data[0]] = {
                "validator_address": data[0],
                "discord_id": campaign_1_discord_id,
                "discord_name": data[2],
                "total_rewards": 0,
                "campaign_1": data[1],
                "referral_code": 0,
                "referral_points": 0,
                "faucet_amount": 0,
                "stake": 0,
                "referrer_discord_id": 0,
                "referrer_discord_name": "",
                "last_received": 0,
                "last_sortition_height": 0,
                "last_bonding_height": 0,
                "unbonding_height": 0,
                "availability_score": 0,
            }


def apply_campaign_1_to_vals(vals_map):
    for campaign_1_discord_id, data in campaign_1.items():
        found = False
        for _, value in vals_map.items():
            campaign_2_discord_id = value["discord_id"]

            if campaign_1_discord_id == campaign_2_discord_id:
                value["campaign_1"] = data[1]
                found = True
                break

        if found == False:
            vals_map[data[0]] = {
                "validator_address": data[0],
                "discord_id": campaign_1_discord_id,
                "discord_name": data[2],
                "total_rewards": 0,
                "campaign_1": data[1],
                "referral_code": 0,
                "referral_points": 0,
                "faucet_amount": 0,
                "stake": 0,
                "referrer_discord_id": 0,
                "referrer_discord_name": "",
                "last_received": 0,
                "last_sortition_height": 0,
                "last_bonding_height": 0,
                "unbonding_height": 0,
                "availability_score": 0,
            }


def apply_referral_points_to_vals(vals_map, referrals):
    for _, ref_data in referrals.items():
        referral_code = ref_data["referral_code"]
        points = ref_data["points"]
        referral_discord_id = ref_data["discord_id"]

        if points == 0:
            continue

        found = False
        for _, value in vals_map.items():
            campaign_2_discord_id = value["discord_id"]

            if referral_discord_id == campaign_2_discord_id:
                value["referral_code"] = referral_code
                value["referral_points"] = points
                found = True
                break

        if found == False:
            print("Something is wrong: " + referral_discord_id)
            # sys.exit(1)


def calculate_rewards_by_user(users_map, vals_map):
    for _, user in users_map.items():
        total_stakes = 0
        total_reward = 0
        total_stakes_online = 0
        total_reward_online = 0

        total_reward = user["campaign_1"] + user["referral_points"]
        total_reward_online = user["campaign_1"] + user["referral_points"]

        for val_addr in user["validators"]:
            val_stake = vals_map[val_addr]["stake"]
            if val_stake == "0":
                print("Something is wrong: " + val_addr)
                sys.exit(1)

            total_stakes += val_stake
            if vals_map[val_addr]["last_received"] > 0:
                if val_addr not in banned_list:
                    total_stakes_online += val_stake
                # else:
                #     print("validator " + val_addr + " is banned.")

        total_reward += total_stakes / 10
        total_reward_online += total_stakes_online / 10

        user["total_stakes"] = total_stakes
        user["total_reward"] = total_reward
        user["total_stakes_online"] = total_stakes_online
        user["total_reward_online"] = total_reward_online


def calculate_rewards_by_val(vals_map):
    for val_addr, val_data in vals_map.items():
        total_rewards = 0

        if val_data["last_received"] > 0:
            if val_addr not in banned_list:
                total_rewards += val_data["stake"] / 10
            # else:
            #     print("validator " + val_addr + " is banned.")

        total_rewards += val_data["campaign_1"]
        total_rewards += val_data["referral_points"]

        val_data["total_rewards"] = total_rewards


if __name__ == "__main__":
    if len(sys.argv) != 4:
        print(
            "Usage: python main.py <validators_path> <referral_path> <network_vals.csv>"
        )
        sys.exit(1)

    validator_path = sys.argv[1]
    referral_path = sys.argv[2]
    network_vals_path = sys.argv[3]

    network_vals = load_csv_as_dict(network_vals_path)
    validator_data = load_json_file(validator_path)
    referral_data = load_json_file(referral_path)

    grpc_validator_channel = grpc.insecure_channel("172.104.46.145:50052")
    grpc_validator_stub = blockchain_pb2_grpc.BlockchainStub(grpc_validator_channel)

    referrals = referrals_by_discord_id(referral_data)
    vals_map = extract_vals_map(validator_data, referrals, network_vals)
    users_map = extract_users_map(validator_data, referrals)

    apply_campaign_1_to_users(users_map)
    apply_campaign_1_to_vals(vals_map)
    apply_referral_points_to_vals(vals_map, referrals)
    calculate_rewards_by_user(users_map, vals_map)
    calculate_rewards_by_val(vals_map)

    write_users_map(users_map)
    write_vals_map(vals_map)
