# mox
Classifying, and physically sorting, Magic: The Gathering cards.

## Things done
- [X] Tools for annotating data
- [X] Tools for exporting data, and doing data augmentation
- [X] Basic inference
- [ ] Accurate classification
  - [X] Card name classification
  - [ ] Set specific classification
- [ ] Inference time optimizations
- [ ] Design of the sorting machine for 3D printing
- [ ] Control of the sorting machine
- [ ] Cataloging capabilities

## Examples
_Card field localisation, before OCR:_
![Localisation example](https://raw.githubusercontent.com/pietroglyph/mox/master/docs/localize_example.png)

## Other documentation
For training your own model, or annotating your own data, see the [applicable readme and files](https://github.com/pietroglyph/mox/tree/master/train).
This should not be necessary for most users; a pretrained inference graph is provided that should give you good results.

## License
All source code in this repository is licensed under the GNU General Public License v3.0, you can get a copy from the filed called `LICENSE` in the top level directory of this repository.

The literal and graphical information contained in this repository about Magic: The Gathering, including card images, the mana symbols, and Oracle text, is copyright Wizards of the Coast, LLC, a subsidiary of Hasbro, Inc. Any content herein is not produced by, endorsed by, supported by, or affiliated with Wizards of the Coast.
