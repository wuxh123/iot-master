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
  NzModalModule, NzPopconfirmModule, NzSelectModule, NzStatisticModule, NzSwitchModule,
  NzTableModule, NzTabsModule,
  NzToolTipModule
} from 'ng-zorro-antd';
import {DashComponent} from './dash/dash.component';
import {MomentModule} from 'ngx-moment';
import {UiModule} from '../ui/ui.module';
import {ChannelComponent} from './channel/channel.component';
import {LinkComponent} from './link/link.component';
import {ChannelEditComponent} from './channel-edit/channel-edit.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {NzSpaceModule} from 'ng-zorro-antd/space';
import {LinkEditComponent} from './link-edit/link-edit.component';
import {PluginComponent} from './plugin/plugin.component';
import {PluginEditComponent} from './plugin-edit/plugin-edit.component';
import {TabsComponent} from './tabs/tabs.component';
import {ModelComponent} from './model/model.component';
import {ModelEditComponent} from './model-edit/model-edit.component';
import {VariableComponent} from './variable/variable.component';
import {VariableEditComponent} from './variable-edit/variable-edit.component';
import {BatchComponent} from './batch/batch.component';
import {BatchEditComponent} from './batch-edit/batch-edit.component';
import {JobComponent} from './job/job.component';
import {JobEditComponent} from './job-edit/job-edit.component';
import {StrategyComponent} from './strategy/strategy.component';
import {StrategyEditComponent} from './strategy-edit/strategy-edit.component';
import {TunnelComponent} from './tunnel/tunnel.component';
import {TunnelEditComponent} from './tunnel-edit/tunnel-edit.component';


@NgModule({
  declarations: [MainComponent, TabsComponent, DashComponent,
    ChannelComponent, ChannelEditComponent,
    LinkComponent, LinkEditComponent, LinkMonitorComponent,
    PluginComponent, PluginEditComponent,
    ModelComponent, ModelEditComponent,
    TunnelComponent, TunnelEditComponent,
    VariableComponent, VariableEditComponent,
    BatchComponent, BatchEditComponent,
    JobComponent, JobEditComponent,
    StrategyComponent, StrategyEditComponent,
  ],
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
    NzStatisticModule,
    NzTabsModule,
  ],
  bootstrap: [MainComponent]
})
export class MainModule {
}
