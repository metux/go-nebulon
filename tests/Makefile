SUBDIRS :=

compile:
	for d in $(SUBDIRS) ; do $(MAKE) -C $$d || exit 1; done

clean:
	rm -Rf *.tmp dirtree/.*.tmp
	for d in $(SUBDIRS) ; do $(MAKE) -C $$d clean || exit 1; done
