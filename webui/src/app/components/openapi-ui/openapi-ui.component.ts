import { OnChanges, Component, ElementRef, Input, ViewChild, AfterViewInit } from '@angular/core';

import SwaggerUI from 'swagger-ui';

@Component({
  selector: 'app-openapi-ui',
  templateUrl: './openapi-ui.component.html',
  styleUrls: []
})
export class OpenAPIUIComponent implements AfterViewInit {

  _spec: Object;
  @ViewChild("openapi") targetdiv: ElementRef;

  constructor(private el: ElementRef) {
  }

  @Input() set spec(val: object) {
    this._spec = val;
    this.apply()
  }

  ngAfterViewInit() {
    this.apply()
  }

  apply() {
    if (!this._spec || !this.targetdiv) {
      return
    }
    // TODO: remove the following line if the import above is fixed
    const SwaggerUI = require('swagger-ui')
    const ui = SwaggerUI({
      spec: this._spec,
      domNode: this.targetdiv.nativeElement,
      deepLinking: true,
      presets: [
        SwaggerUI.presets.apis
      ],
    });
  }
}
