import { Component, OnInit, Output, EventEmitter, Input, ElementRef, OnDestroy, Renderer2 } from '@angular/core';
import { DomHandler } from 'primeng/primeng';
import { HelpsStoreService } from '../../stores/help-store.service';

@Component({
  selector: 'app-help-widget',
  templateUrl: './help-widget.component.html',
  styleUrls: ['./help-widget.component.css'],
  providers: [DomHandler]
})
export class HelpWidgetComponent implements OnInit {

  private _key = '';
  @Output() select: EventEmitter<any> = new EventEmitter();

  constructor(
    public help: HelpsStoreService
  ) {
  }

  ngOnInit() {
  }

  public get key() {
    return this._key;
  }

  @Input() public set key(val: string) {
    this._key = val;
  }
}
