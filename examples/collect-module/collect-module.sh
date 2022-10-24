# collect theme as module
start to create `mytheme` module
# project module has no owner with default mounts
&main.moduleAdapter{
  projectMod:true, owner:main.Module(nil),
  mounts:[]main.Mount{
    main.Mount{
      Source:"mycontent", Target:"content",Lang:"en"},
    main.Mount{Source:"data", Target:"data", Lang:""},
    main.Mount{
      Source:"layouts", Target:"layouts", Lang:""},
    main.Mount{Source:"i18n", Target:"i18n", Lang:""},
    main.Mount{
      Source:"archetypes", Target:"archetypes", Lang:""},
    main.Mount{Source:"assets", Target:"assets", Lang:""},
    main.Mount{
      Source:"static", Target:"static", Lang:""}},
    config:main.Config{Mounts:[]main.Mount(nil),
    Imports:[]main.Import{main.Import{Path:"mytheme"}}}}
# theme module owned by project module
# with no import in the example
&main.moduleAdapter{projectMod:false,
owner:(*main.moduleAdapter)(0xc000102120),
mounts:[]main.Mount(nil),
config:main.Config{
  Mounts:[]main.Mount(nil), Imports:[]main.Import(nil)}}

Program exited.