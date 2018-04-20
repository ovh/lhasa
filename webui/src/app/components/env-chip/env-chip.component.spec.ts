import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EnvChipComponent } from './env-chip.component';

describe('EnvChipComponent', () => {
  let component: EnvChipComponent;
  let fixture: ComponentFixture<EnvChipComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EnvChipComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EnvChipComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
