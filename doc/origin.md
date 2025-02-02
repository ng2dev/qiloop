## Apendice: analysis with Wireshark

To analyse the payload, one can capture the QiMessaging packets with
tcpdump or wireshark like this:

    $ tcpdump -i lo -w qicli-info.pcap port 9559

And analyse them with a plugin for Wireshark which to decode headers of the messages.
The plugin must be installed into ``$HOME/.config/wireshark/plugins/``:

    https://github.com/aldebaran/libqi/tree/team/platform/dev/tools/wireshark



### 1. Call: ServiceServer.Authenticate

    00000000: 0400 0000 1200 0000 436c 6965 6e74 5365  ........ClientSe
    00000010: 7276 6572 536f 636b 6574 0100 0000 6201  rverSocket....b.
    00000020: 0c00 0000 4d65 7373 6167 6546 6c61 6773  ....MessageFlags
    00000030: 0100 0000 6201 0f00 0000 4d65 7461 4f62  ....b.....MetaOb
    00000040: 6a65 6374 4361 6368 6501 0000 0062 0115  jectCache....b..
    00000050: 0000 0052 656d 6f74 6543 616e 6365 6c61  ...RemoteCancela
    00000060: 626c 6543 616c 6c73 0100 0000 6201       bleCalls....b.

    Strings:

        int: 4 // number of elements

        int: 18 // size of the string
        string: ClientServerSocket // element

        int: 1 // size (number of signature)
        char: 'b' // signature
        char: 0x1 // value: true

        int: 12
        string: MessageFlags
        int: 1
        char: 'b'
        char: 0x1

        int: 15
        string: MetaObjectCache
        int: 1
        char: 'b'
        char: 0x1

        int: 18
        string: RemoteCancelableCalls
        int: 1
        char: 'b'
        char: 0x1

    std::map<std::string, AnyValue> m = {
        {"ClientServerSocket", "boolean:true"},
        {"MessageFlags", "boolean:true"},
        {"MetaObjectCache", "boolean:true"},
        {"RemoteCancelableCalls", "boolean:true"}
    };


### 2. Reply: ServiceServer.Authenticate

    00000000: 0500 0000 1200 0000 436c 6965 6e74 5365  ........ClientSe
    00000010: 7276 6572 536f 636b 6574 0100 0000 6201  rverSocket....b.
    00000020: 0c00 0000 4d65 7373 6167 6546 6c61 6773  ....MessageFlags
    00000030: 0100 0000 6201 0f00 0000 4d65 7461 4f62  ....b.....MetaOb
    00000040: 6a65 6374 4361 6368 6501 0000 0062 0115  jectCache....b..
    00000050: 0000 0052 656d 6f74 6543 616e 6365 6c61  ...RemoteCancela
    00000060: 626c 6543 616c 6c73 0100 0000 6201 0f00  bleCalls....b...
    00000070: 0000 5f5f 7169 5f61 7574 685f 7374 6174  ..__qi_auth_stat
    00000080: 6501 0000 0069 0300 0000                 e....i....

    Strings:

        int: 5 // number of elements

        int: 18 // size of the string
        string: ClientServerSocket // element
        int: 1 // size (number of signature)
        char: 'b' // signature
        char: 0x1 // value: true

        int: 12
        string: MessageFlags
        int: 1
        char: 'b'
        char: 0x1

        int: 15
        string: MetaObjectCache
        int: 1
        char: 'b'
        char: 0x1

        int: 18
        string: RemoteCancelableCalls
        int: 1
        char: 'b'
        char: 0x1

        int: 15
        string: __qi_auth_state
        int: 1
        char: 'i'
        int: 0x3


    std::map<std::string, AnyValue> m = {
        {"ClientServerSocket", "boolean:true"},
        {"MessageFlags", "boolean:true"},
        {"MetaObjectCache", "boolean:true"},
        {"RemoteCancelableCalls", "boolean:true"},
        {"__qi_auth_state", "int:3"}
    };

### 3.Call: ServiceDirectory.metaObject(int): MetaObject


### 4.Reply: ServiceDirectory.metaObject(int): MetaObject

    00000000: 1600 0000 0000 0000 0000 0000 0100 0000  ................

    int: 22 // number of elements in the map
    int: 0 // index 0 of the map
    int: 0 // MetaMethod.uuid
    int: 1 // size of MetaMethod.returnSignature

    00000010: 4c0d 0000 0072 6567 6973 7465 7245 7665  L....registerEve

    char: 'L' // MetaMethod.returnSignature
    int: 13 // size of MetaMethod.name
    string: registerEvent // MetaMethod.name

    00000020: 6e74 0500 0000 2849 494c 2900 0000 0000  nt....(IIL).....

    int: 5 // size of MetaMethod.parametersSignature
    string: "(IIL)" // MetaMethod.parametersSignature
    int: 0 // size of description
    int: 0 // size of MetaMethod.MetaMethodParameter.name

    <0><0><1>"L"<13>"registerEvent"<5>"(IIL)"

    00000030: 0000 0000 0000 0001 0000 0001 0000 0001  ................

    int: 0 // size of MetaMethod.MetaMethodParameter.description
    int: 1 // index 1 of the methods map
    int: 1 // MetaMethod.uuid
    int: 1 // MetaMethod.returnSignature

    00000040: 0000 0076 0f00 0000 756e 7265 6769 7374  ...v....unregist

    char: 'V'
    int: 15 // size of MetaMethod.name
    string: "unregisterEvent"

    00000050: 6572 4576 656e 7405 0000 0028 4949 4c29  erEvent....(IIL)

    int: 5 // size of MetaMethod.parametersSignature
    string: "(IIL)"

    <0><1><1><1>"V"<15>"unregisterEvent"<5>"(IIL)"

    00000060: 0000 0000 0000 0000 0000 0000 0200 0000  ................

    int: 0 // size of MetaMethod.description
    int: 0 // size of MetaMethod.Parameters[]
    int: 0 // size of MetaMethod.ReturnDescription
    int: 2 // index 2?

    00000070: 0200 0000 1f01 0000 287b 4928 4973 7373  ........({I(Isss

    int: 2 // MetaMethod.uuid
    int: 287 // size of string
    string: "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>"


    "(ss)<MetaMethodParameter,name,description>" :
    struct MetaMethodParameter {
        string name;
        string description;
    }
    "{(Issss[MetaMethodParameter...]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}"
    struct MetaMethod {
        int uuid;
        string returnSignature;
        string name;
        string parametersSignature;
        string description;
        MetaMethodParameter parameters;
        string returnDescription;
    }
    "{I(Iss)<MetaSignal,uid,name,signature>}"

    struct MetaSignal {
        int uuid;
        string name;
        string signature;
    }
    "{I(Iss)<MetaProperty,uid,name,signature>}"
    struct MetaProperty {
        int uuid;
        string name;
        string signature;
    }
    "({ ... MetaMethod ... }{ ... MetaSignal ... }{ ... MetaProperty ...}s)<MetaObject,methods,signals,properties,description>"

    struct MetaObject {
        Map[Int,MetaMethod] methods
        Map[Int,MetaSignal] signals
        Map[Int,MetaProperty] properties
        string description
    }


    00000080: 735b 2873 7329 3c4d 6574 614d 6574 686f  s[(ss)<MetaMetho
    00000090: 6450 6172 616d 6574 6572 2c6e 616d 652c  dParameter,name,
    000000a0: 6465 7363 7269 7074 696f 6e3e 5d73 293c  description>]s)<
    000000b0: 4d65 7461 4d65 7468 6f64 2c75 6964 2c72  MetaMethod,uid,r
    000000c0: 6574 7572 6e53 6967 6e61 7475 7265 2c6e  eturnSignature,n
    000000d0: 616d 652c 7061 7261 6d65 7465 7273 5369  ame,parametersSi
    000000e0: 676e 6174 7572 652c 6465 7363 7269 7074  gnature,descript
    000000f0: 696f 6e2c 7061 7261 6d65 7465 7273 2c72  ion,parameters,r
    00000100: 6574 7572 6e44 6573 6372 6970 7469 6f6e  eturnDescription
    00000110: 3e7d 7b49 2849 7373 293c 4d65 7461 5369  >}{I(Iss)<MetaSi
    00000120: 676e 616c 2c75 6964 2c6e 616d 652c 7369  gnal,uid,name,si
    00000130: 676e 6174 7572 653e 7d7b 4928 4973 7329  gnature>}{I(Iss)
    00000140: 3c4d 6574 6150 726f 7065 7274 792c 7569  <MetaProperty,ui
    00000150: 642c 6e61 6d65 2c73 6967 6e61 7475 7265  d,name,signature
    00000160: 3e7d 7329 3c4d 6574 614f 626a 6563 742c  >}s)<MetaObject,
    00000170: 6d65 7468 6f64 732c 7369 676e 616c 732c  methods,signals,
    00000180: 7072 6f70 6572 7469 6573 2c64 6573 6372  properties,descr
    00000190: 6970 7469 6f6e 3e0a 0000 006d 6574 614f  iption>....metaO

    int: 10 // size of string
    string: "metaObject" // name of the method

    000001a0: 626a 6563 7403 0000 0028 4929 0000 0000  bject....(I)....

    int: 3 // size of string // parametersSignature
    string: "(I)" // unsigned int
    unsigned int: 0 // size of description

    <0><0><0><2><2><287>"{(Issss[Meta..."<10>"metaObject"<3>"(I)"

    000001b0: 0000 0000 0000 0000 0300 0000 0300 0000  ................

    int: 0 // size of Parameters
    int: 0 // size of ReturnDescription
    int: 3 // index of map
    int: 3 // MetaMethod.uuid

    000001c0: 0100 0000 7609 0000 0074 6572 6d69 6e61  ....v....termina

    int: 1 // size of string
    string: "v" // void
    int: 9 // size of string
    string: "terminate"

    000001d0: 7465 0300 0000 2849 2900 0000 0000 0000  te....(I).......

    int: 3 // size of string
    string: "(I)" // unsigned int
    int: 0
    int: 0

    000001e0: 0000 0000 0005 0000 0005 0000 0001 0000  ................

    int: 0
    int: 5
    int: 5
    int: 1 // size of string

    000001f0: 006d 0800 0000 7072 6f70 6572 7479 0300  .m....property..

    string: "m"
    int: 8 // size of string
    sting: "property"
    int: 3 // size of string

    00000200: 0000 286d 2900 0000 0000 0000 0000 0000  ..(m)...........

    sting: "(m)"
    int: 0
    int: 0
    int: 0

    00000210: 0006 0000 0006 0000 0001 0000 0076 0b00  .............v..

    int: 6
    int: 6
    int: 1 // size of string
    string: "v"
    int: 11 // size of string

    00000220: 0000 7365 7450 726f 7065 7274 7904 0000  ..setProperty...

    string: "setProperty"
    int: 4 // size of string

    00000230: 0028 6d6d 2900 0000 0000 0000 0000 0000  .(mm)...........
    00000240: 0007 0000 0007 0000 0003 0000 005b 735d  .............[s]
    00000250: 0a00 0000 7072 6f70 6572 7469 6573 0200  ....properties..
    00000260: 0000 2829 0000 0000 0000 0000 0000 0000  ..()............
    00000270: 0800 0000 0800 0000 0100 0000 4c1a 0000  ............L...
    00000280: 0072 6567 6973 7465 7245 7665 6e74 5769  .registerEventWi
    00000290: 7468 5369 676e 6174 7572 6506 0000 0028  thSignature....(
    000002a0: 4949 4c73 2900 0000 0000 0000 0000 0000  IILs)...........
    000002b0: 0050 0000 0050 0000 0001 0000 0062 0e00  .P...P.......b..
    000002c0: 0000 6973 5374 6174 7345 6e61 626c 6564  ..isStatsEnabled
    000002d0: 0200 0000 2829 0000 0000 0000 0000 0000  ....()..........
    000002e0: 0000 5100 0000 5100 0000 0100 0000 760b  ..Q...Q.......v.
    000002f0: 0000 0065 6e61 626c 6553 7461 7473 0300  ...enableStats..
    00000300: 0000 2862 2900 0000 0000 0000 0000 0000  ..(b)...........
    00000310: 0052 0000 0052 0000 00c2 0000 007b 4928  .R...R.......{I(
    00000320: 4928 6666 6629 3c4d 696e 4d61 7853 756d  I(fff)<MinMaxSum
    00000330: 2c6d 696e 5661 6c75 652c 6d61 7856 616c  ,minValue,maxVal
    00000340: 7565 2c63 756d 756c 6174 6564 5661 6c75  ue,cumulatedValu
    00000350: 653e 2866 6666 293c 4d69 6e4d 6178 5375  e>(fff)<MinMaxSu
    00000360: 6d2c 6d69 6e56 616c 7565 2c6d 6178 5661  m,minValue,maxVa
    00000370: 6c75 652c 6375 6d75 6c61 7465 6456 616c  lue,cumulatedVal
    00000380: 7565 3e28 6666 6629 3c4d 696e 4d61 7853  ue>(fff)<MinMaxS
    00000390: 756d 2c6d 696e 5661 6c75 652c 6d61 7856  um,minValue,maxV
    000003a0: 616c 7565 2c63 756d 756c 6174 6564 5661  alue,cumulatedVa
    000003b0: 6c75 653e 293c 4d65 7468 6f64 5374 6174  lue>)<MethodStat
    000003c0: 6973 7469 6373 2c63 6f75 6e74 2c77 616c  istics,count,wal
    000003d0: 6c2c 7573 6572 2c73 7973 7465 6d3e 7d05  l,user,system>}.
    000003e0: 0000 0073 7461 7473 0200 0000 2829 0000  ...stats....()..
    000003f0: 0000 0000 0000 0000 0000 5300 0000 5300  ..........S...S.
    00000400: 0000 0100 0000 760a 0000 0063 6c65 6172  ......v....clear
    00000410: 5374 6174 7302 0000 0028 2900 0000 0000  Stats....().....
    00000420: 0000 0000 0000 0054 0000 0054 0000 0001  .......T...T....
    00000430: 0000 0062 0e00 0000 6973 5472 6163 6545  ...b....isTraceE
    00000440: 6e61 626c 6564 0200 0000 2829 0000 0000  nabled....()....
    00000450: 0000 0000 0000 0000 5500 0000 5500 0000  ........U...U...
    00000460: 0100 0000 760b 0000 0065 6e61 626c 6554  ....v....enableT
    00000470: 7261 6365 0300 0000 2862 2900 0000 0000  race....(b).....
    00000480: 0000 0000 0000 0064 0000 0064 0000 004e  .......d...d...N
    00000490: 0000 0028 7349 7349 5b73 5d73 293c 5365  ...(sIsI[s]s)<Se
    000004a0: 7276 6963 6549 6e66 6f2c 6e61 6d65 2c73  rviceInfo,name,s
    000004b0: 6572 7669 6365 4964 2c6d 6163 6869 6e65  erviceId,machine
    000004c0: 4964 2c70 726f 6365 7373 4964 2c65 6e64  Id,processId,end
    000004d0: 706f 696e 7473 2c73 6573 7369 6f6e 4964  points,sessionId
    000004e0: 3e07 0000 0073 6572 7669 6365 0300 0000  >....service....
    000004f0: 2873 2900 0000 0000 0000 0000 0000 0065  (s)............e
    00000500: 0000 0065 0000 0050 0000 005b 2873 4973  ...e...P...[(sIs
    00000510: 495b 735d 7329 3c53 6572 7669 6365 496e  I[s]s)<ServiceIn
    00000520: 666f 2c6e 616d 652c 7365 7276 6963 6549  fo,name,serviceI
    00000530: 642c 6d61 6368 696e 6549 642c 7072 6f63  d,machineId,proc
    00000540: 6573 7349 642c 656e 6470 6f69 6e74 732c  essId,endpoints,
    00000550: 7365 7373 696f 6e49 643e 5d08 0000 0073  sessionId>]....s
    00000560: 6572 7669 6365 7302 0000 0028 2900 0000  ervices....()...
    00000570: 0000 0000 0000 0000 0066 0000 0066 0000  .........f...f..
    00000580: 0001 0000 0049 0f00 0000 7265 6769 7374  .....I....regist
    00000590: 6572 5365 7276 6963 6550 0000 0028 2873  erServiceP...((s
    000005a0: 4973 495b 735d 7329 3c53 6572 7669 6365  IsI[s]s)<Service
    000005b0: 496e 666f 2c6e 616d 652c 7365 7276 6963  Info,name,servic
    000005c0: 6549 642c 6d61 6368 696e 6549 642c 7072  eId,machineId,pr
    000005d0: 6f63 6573 7349 642c 656e 6470 6f69 6e74  ocessId,endpoint
    000005e0: 732c 7365 7373 696f 6e49 643e 2900 0000  s,sessionId>)...
    000005f0: 0000 0000 0000 0000 0067 0000 0067 0000  .........g...g..
    00000600: 0001 0000 0076 1100 0000 756e 7265 6769  .....v....unregi
    00000610: 7374 6572 5365 7276 6963 6503 0000 0028  sterService....(
    00000620: 4929 0000 0000 0000 0000 0000 0000 6800  I)............h.
    00000630: 0000 6800 0000 0100 0000 760c 0000 0073  ..h.......v....s
    00000640: 6572 7669 6365 5265 6164 7903 0000 0028  erviceReady....(
    00000650: 4929 0000 0000 0000 0000 0000 0000 6900  I)............i.
    00000660: 0000 6900 0000 0100 0000 7611 0000 0075  ..i.......v....u
    00000670: 7064 6174 6553 6572 7669 6365 496e 666f  pdateServiceInfo
    00000680: 5000 0000 2828 7349 7349 5b73 5d73 293c  P...((sIsI[s]s)<
    00000690: 5365 7276 6963 6549 6e66 6f2c 6e61 6d65  ServiceInfo,name
    000006a0: 2c73 6572 7669 6365 4964 2c6d 6163 6869  ,serviceId,machi
    000006b0: 6e65 4964 2c70 726f 6365 7373 4964 2c65  neId,processId,e
    000006c0: 6e64 706f 696e 7473 2c73 6573 7369 6f6e  ndpoints,session
    000006d0: 4964 3e29 0000 0000 0000 0000 0000 0000  Id>)............
    000006e0: 6c00 0000 6c00 0000 0100 0000 7309 0000  l...l.......s...
    000006f0: 006d 6163 6869 6e65 4964 0200 0000 2829  .machineId....()
    00000700: 0000 0000 0000 0000 0000 0000 6d00 0000  ............m...
    00000710: 6d00 0000 0100 0000 6f10 0000 005f 736f  m.......o...._so
    00000720: 636b 6574 4f66 5365 7276 6963 6503 0000  cketOfService...
    00000730: 0028 4929 0000 0000 0000 0000 0000 0000  .(I)............
    00000740: 0300 0000 5600 0000 5600 0000 0b00 0000  ....V...V.......
    00000750: 7472 6163 654f 626a 6563 748b 0000 0028  traceObject....(
    00000760: 2849 6949 6d28 6c6c 293c 7469 6d65 7661  (IiIm(ll)<timeva
    00000770: 6c2c 7476 5f73 6563 2c74 765f 7573 6563  l,tv_sec,tv_usec
    00000780: 3e6c 6c49 4929 3c45 7665 6e74 5472 6163  >llII)<EventTrac
    00000790: 652c 6964 2c6b 696e 642c 736c 6f74 4964  e,id,kind,slotId
    000007a0: 2c61 7267 756d 656e 7473 2c74 696d 6573  ,arguments,times
    000007b0: 7461 6d70 2c75 7365 7255 7354 696d 652c  tamp,userUsTime,
    000007c0: 7379 7374 656d 5573 5469 6d65 2c63 616c  systemUsTime,cal
    000007d0: 6c65 7243 6f6e 7465 7874 2c63 616c 6c65  lerContext,calle
    000007e0: 6543 6f6e 7465 7874 3e29 6a00 0000 6a00  eContext>)j...j.
    000007f0: 0000 0c00 0000 7365 7276 6963 6541 6464  ......serviceAdd
    00000800: 6564 0400 0000 2849 7329 6b00 0000 6b00  ed....(Is)k...k.
    00000810: 0000 0e00 0000 7365 7276 6963 6552 656d  ......serviceRem
    00000820: 6f76 6564 0400 0000 2849 7329 0000 0000  oved....(Is)....
    00000830: 0000 0000 0a                             .....

