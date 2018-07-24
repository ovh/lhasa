import { OnChanges, Component, ElementRef, Input, ViewChild, AfterViewInit } from '@angular/core';

import SwaggerUI from 'swagger-ui';

@Component({
  selector: 'app-openapi-ui',
  templateUrl: './openapi-ui.component.html',
  styleUrls: []
})
export class OpenAPIUIComponent implements AfterViewInit {

  _url: Object;
  @ViewChild('openapi') targetdiv: ElementRef;

  constructor(private el: ElementRef) {
  }

  @Input() set url(val: string) {
    this._url = val;
    this.apply();
  }

  ngAfterViewInit() {
    this.apply();
  }

  apply() {
    if (!this._url || !this.targetdiv) {
      return;
    }
    // TODO: remove the following line if the import above is fixed
    const SwaggerUI = require('swagger-ui');
    const ui = SwaggerUI({
      url: this._url,
      domNode: this.targetdiv.nativeElement,
      deepLinking: true,
      validatorUrl: null,
      presets: [
        SwaggerUI.presets.apis
      ],
    });
  }
}
