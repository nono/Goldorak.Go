include $(GOROOT)/src/Make.$(GOARCH)

TARG=goldorak
GOFILES=\
		action.go\
		config.go\
		goldorak.go\
		helpers.go\
		model.go\

include $(GOROOT)/src/Make.pkg
