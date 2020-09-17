import {NgModule} from '@angular/core';

import {IconsProviderModule} from './icons-provider.module';
import {NzLayoutModule} from 'ng-zorro-antd/layout';
import {NzMenuModule} from 'ng-zorro-antd/menu';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {CommonModule} from '@angular/common';
import {HttpClientModule} from '@angular/common/http';
import {MainRoutingModule} from './main-routing.module';

import {MainComponent} from './main.component';
import {
  NzButtonModule,
  NzCheckboxModule, NzDividerModule, NzDrawerModule,
  NzFormModule,
  NzIconModule,
  NzInputModule, NzInputNumberModule,
  NzModalModule, NzPopconfirmModule, NzSelectModule, NzSwitchModule,
  NzTableModule,
  NzToolTipModule
} from 'ng-zorro-antd';
import {DashComponent} from './dash/dash.component';
import {MomentModule} from 'ngx-moment';
import {UiModule} from '../ui/ui.module';
import {ChannelComponent} from './channel/channel.component';
import {LinkComponent} from './link/link.component';
import {ChannelEditComponent} from './channel-edit/channel-edit.component';
import {ChannelDetailComponent} from './channel-detail/channel-detail.component';
import {LinkDetailComponent} from './link-detail/link-detail.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {NzSpaceModule} from 'ng-zorro-antd/space';
import {LinkEditComponent} from './link-edit/link-edit.component';
import {PluginComponent} from './plugin/plugin.component';
// import {MqttModule} from "ngx-mqtt";
// import {IMqttServiceOptions} from "ngx-mqtt/lib/mqtt.model";
//
// export const MQTT_SERVICE_OPTIONS: IMqttServiceOptions = {
//   url: 'ws://localhost:8080/api/mqtt'
// };


@NgModule({
  declarations: [MainComponent, DashComponent,
    ChannelComponent, ChannelDetailComponent, ChannelEditComponent,
    LinkComponent, LinkDetailComponent, LinkMonitorComponent, LinkEditComponent,
    PluginComponent],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    MomentModule,
    UiModule,
    // IconsProviderModule,
    // BrowserModule,
    NzIconModule,
    NzLayoutModule,
    NzMenuModule,
    HttpClientModule,
    MainRoutingModule,
    NzToolTipModule,
    NzTableModule,
    NzModalModule,
    NzFormModule,
    NzButtonModule,
    NzInputModule,
    NzCheckboxModule,
    NzSwitchModule,
    NzPopconfirmModule,
    IconsProviderModule,
    NzDividerModule,
    NzDrawerModule,
    NzSelectModule,
    NzSpaceModule,
    NzInputNumberModule,
    //MqttModule.forRoot(MQTT_SERVICE_OPTIONS)
  ],
  bootstrap: [MainComponent]
})
export class MainModule {
}
