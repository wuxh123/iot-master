import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ToolbarComponent} from './toolbar/toolbar.component';
import { FormGridComponent } from './form-grid/form-grid.component';
import {NzGridModule} from 'ng-zorro-antd';
import { FormItemComponent } from './form-item/form-item.component';


@NgModule({
  declarations: [
    ToolbarComponent,
    FormGridComponent,
    FormItemComponent,
  ],
  exports: [
    ToolbarComponent,
    FormGridComponent,
    FormItemComponent,
  ],
  imports: [
    CommonModule,
    // NzLayoutModule,
    NzGridModule,
  ]
})
export class UiModule {
}
