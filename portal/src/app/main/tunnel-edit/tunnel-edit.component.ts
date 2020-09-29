import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzDrawerRef} from 'ng-zorro-antd';

@Component({
  selector: 'app-tunnel-edit',
  templateUrl: './tunnel-edit.component.html',
  styleUrls: ['./tunnel-edit.component.scss']
})
export class TunnelEditComponent implements OnInit {

  @Input() tunnel: any = {};


  constructor(private as: ApiService, private drawerRef: NzDrawerRef<string>) {
  }

  ngOnInit(): void {
    if (this.tunnel.id) {
      this.as.get('tunnel/' + this.tunnel.id).subscribe(res => {
        this.tunnel = res.data;
      });
    }
  }

  submit(): void {
    if (this.tunnel.id) {
      this.as.put('tunnel/' + this.tunnel.id, this.tunnel).subscribe(res => {
        console.log(res);
        // TODO 修改成功
        this.drawerRef.close(res.data);
      });
    } else {
      this.as.post('tunnel', this.tunnel).subscribe(res => {
        console.log(res);
        // TODO 保存成功
        this.drawerRef.close(res.data);
      });
    }
  }
}
