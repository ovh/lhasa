import { OnChanges, Component, ElementRef, Input, ViewChild, AfterViewInit } from '@angular/core';

import SwaggerUI from 'swagger-ui';

@Component({
  selector: 'app-openapi-ui',
  templateUrl: './openapi-ui.component.html',
  styleUrls: []
})
export class OpenAPIUIComponent implements AfterViewInit {

  _url: string;
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
    SwaggerUI({
      url: this._url,
      domNode: this.targetdiv.nativeElement,
      deepLinking: false,
      validatorUrl: null,
      displayRequestDuration: true,
      presets: [
        SwaggerUI.presets.apis
      ],
    });
  }
}
